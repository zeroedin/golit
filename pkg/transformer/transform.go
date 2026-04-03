// Package transformer provides file-level and directory-level HTML
// transformation, walking HTML files on disk and expanding custom
// elements using the QJS rendering engine.
package transformer

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/sspriggs/golit/pkg/fileutil"
	"github.com/sspriggs/golit/pkg/jsengine"
)

// Options configures the transformer.
type Options struct {
	// DefsDir is the directory containing .golit.bundle.js files (Mode 1).
	DefsDir string

	// CompiledFile is a path to a single .golit.compiled.js artifact
	// containing all bundles and a tag registry manifest.
	CompiledFile string

	// SourcesDir is a directory of component .js/.ts source files to bundle
	// on-demand (Mode 2).
	SourcesDir string

	// ImportMapFile is a path to an import map JSON file (Mode 3).
	ImportMapFile string

	// AutoDiscover enables HTML auto-discovery of <script type="importmap">
	// and <script type="module"> tags (Mode 4). Enabled by default when
	// no other discovery mode is specified.
	AutoDiscover bool

	// Ignored is a set of custom element tag names to skip during SSR.
	// These elements will be left as-is for client-side rendering only.
	Ignored map[string]bool

	// Preload is a list of extra JS modules to bundle and load into the
	// QJS engine before component rendering. Each entry is a bare module
	// specifier (e.g. "prism-esm") or a file path.
	Preload []string

	// Verbose prints progress information to stderr.
	Verbose bool

	// DryRun reads and transforms files but does not write them back.
	DryRun bool

	// Concurrency is the number of parallel workers for file processing.
	// 0 or 1 means sequential (default). Set to runtime.NumCPU() or
	// a specific value for parallel processing of large sites.
	Concurrency int

	// OutDir is an optional output directory. When set, transformed files
	// are written here instead of modifying the input files in-place.
	OutDir string

	// Isolate creates a fresh QJS context per HTML file, clearing all
	// global state between files. Slower but safer for untrusted components.
	Isolate bool
}

// RenderError records a custom element that failed to render during SSR.
type RenderError struct {
	TagName string
	File    string
	Err     error
}

func (e RenderError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("<%s> in %s: %v", e.TagName, e.File, e.Err)
	}
	return fmt.Sprintf("<%s>: %v", e.TagName, e.Err)
}

// Result holds stats from a transform run.
type Result struct {
	FilesProcessed int
	FilesModified  int
	Errors         []error
	RenderErrors   []RenderError
	Unregistered   []string
}

// TransformDir processes all HTML files in a directory tree.
//
// Processing happens in two passes:
//  1. Discovery (sequential): read every HTML file and run component
//     auto-discovery so the registry is fully populated.
//  2. Render (parallel when Concurrency > 1): transform files using a
//     pool of QJS engines, one engine per worker goroutine.
func TransformDir(dir string, opts Options) (*Result, error) {
	registry := jsengine.NewRegistry()

	// Mode 0: Pre-compiled single artifact
	if opts.CompiledFile != "" {
		if err := registry.LoadCompiled(opts.CompiledFile); err != nil {
			return nil, fmt.Errorf("loading compiled artifact: %w", err)
		}
	}

	// Mode 1: Pre-bundled .golit.bundle.js files
	if opts.DefsDir != "" {
		if err := registry.LoadDir(opts.DefsDir); err != nil {
			return nil, fmt.Errorf("loading bundles: %w", err)
		}
	}

	// Mode 2: Source directory -- bundle all .js/.ts files on-demand
	if opts.SourcesDir != "" {
		if err := registry.LoadSourceDir(opts.SourcesDir); err != nil {
			return nil, fmt.Errorf("loading sources: %w", err)
		}
	}

	// Mode 3: CLI import map -- will be used with HTML discovery
	var cliImportMap *jsengine.ImportMap
	if opts.ImportMapFile != "" {
		im, err := jsengine.LoadImportMapFile(opts.ImportMapFile)
		if err != nil {
			return nil, fmt.Errorf("loading import map: %w", err)
		}
		cliImportMap = im
	}

	// Auto-discover is on by default if no other mode specified
	autoDiscover := opts.AutoDiscover
	if !autoDiscover && opts.DefsDir == "" && opts.CompiledFile == "" && opts.SourcesDir == "" && opts.ImportMapFile == "" {
		autoDiscover = true
	}

	// Collect HTML files
	htmlFiles, err := collectHTMLFiles(dir)
	if err != nil {
		return nil, fmt.Errorf("collecting HTML files: %w", err)
	}

	// Resolve and bundle preload modules (shared across all engines).
	var preloadBundles []string
	for _, mod := range opts.Preload {
		modPath, err := jsengine.ResolveModulePath(mod, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: could not resolve preload %s: %v\n", mod, err)
			continue
		}
		bundle, err := jsengine.BundlePreload(modPath, mod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: could not bundle preload %s: %v\n", mod, err)
			continue
		}
		preloadBundles = append(preloadBundles, bundle)
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "golit: preloaded %s from %s\n", mod, modPath)
		}
	}

	// ── Pass 1: Discovery (sequential) ──────────────────────────────
	// Read every HTML file and run component auto-discovery so the
	// registry is fully populated before rendering begins.
	if autoDiscover || cliImportMap != nil {
		for _, filePath := range htmlFiles {
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue // will be reported in the render pass
			}
			htmlDir := filepath.Dir(filePath)
			discoverFromHTML(string(data), htmlDir, dir, registry, cliImportMap, opts.Verbose)
		}
	}

	// ── Pass 2: Render ──────────────────────────────────────────────
	// Determine concurrency. Isolate mode forces sequential processing
	// because each file needs a fresh engine.
	// Concurrency: 1 = sequential (default), N>1 = N parallel workers.
	// Isolate mode forces sequential since each file needs a fresh engine.
	workers := opts.Concurrency
	if workers < 1 {
		workers = 1
	}
	if opts.Isolate {
		workers = 1
	}
	// No point having more workers than files.
	if workers > len(htmlFiles) {
		workers = len(htmlFiles)
	}
	if workers < 1 {
		workers = 1
	}

	if opts.Verbose && workers > 1 {
		fmt.Fprintf(os.Stderr, "golit: using %d parallel workers\n", workers)
	}

	var result *Result
	if workers == 1 {
		result, err = transformSequential(htmlFiles, dir, registry, opts, preloadBundles)
	} else {
		result, err = transformParallel(htmlFiles, dir, registry, opts, preloadBundles, workers)
	}
	if err != nil {
		return nil, err
	}

	// Filter the unregistered list: remove sub-components of known elements.
	// We need a single engine to check QJS registration for side-effect-defined
	// elements (e.g. sub-components registered by a parent bundle).
	checkEngine, engineErr := jsengine.NewEngine()
	if engineErr == nil {
		defer checkEngine.Close()
		checkEngine.SetPreloadModules(opts.Preload)
		for _, pb := range preloadBundles {
			_ = checkEngine.LoadBundle(pb)
		}
		for _, tag := range registry.TagNames() {
			checkEngine.LoadBundleForTag(tag, registry)
		}
	}
	knownTags := registry.TagNames()
	knownPaths := registry.ProcessedPaths()
	var trulyUnregistered []string
	for _, tag := range registry.Unregistered() {
		if checkEngine != nil && checkEngine.IsRegistered(tag) {
			continue
		}

		// Check prefix match against known tags
		isSub := false
		for _, known := range knownTags {
			if strings.HasPrefix(tag, known+"-") {
				isSub = true
				break
			}
		}
		if isSub {
			continue
		}

		// Check if tag's .js file exists in any known element's directory
		for _, knownPath := range knownPaths {
			siblingPath := filepath.Join(filepath.Dir(knownPath), tag+".js")
			if _, err := os.Stat(siblingPath); err == nil {
				isSub = true
				break
			}
		}
		if isSub {
			continue
		}

		trulyUnregistered = append(trulyUnregistered, tag)
	}
	result.Unregistered = trulyUnregistered

	return result, nil
}

// initEngine creates an engine and loads preload bundles into it.
func initEngine(preloadBundles []string, preloadModules []string) (*jsengine.Engine, error) {
	engine, err := jsengine.NewEngine()
	if err != nil {
		return nil, err
	}
	engine.SetPreloadModules(preloadModules)
	for _, pb := range preloadBundles {
		if err := engine.LoadBundle(pb); err != nil {
			engine.Close()
			return nil, err
		}
	}
	return engine, nil
}

// transformSequential processes files one at a time with a single engine.
// Used when Concurrency==1 or Isolate mode is on.
func transformSequential(htmlFiles []string, dir string, registry *jsengine.Registry, opts Options, preloadBundles []string) (*Result, error) {
	engine, err := initEngine(preloadBundles, opts.Preload)
	if err != nil {
		return nil, fmt.Errorf("creating JS engine: %w", err)
	}
	defer engine.Close()

	var (
		processed    int
		modified     int
		errorsList   []error
		renderErrors []RenderError
	)

	for _, filePath := range htmlFiles {
		if opts.Isolate {
			if err := engine.Reset(); err != nil {
				errorsList = append(errorsList, fmt.Errorf("resetting engine for %s: %w", filePath, err))
				processed++
				continue
			}
			// Re-load preloads after reset.
			engine.SetPreloadModules(opts.Preload)
			for _, pb := range preloadBundles {
				_ = engine.LoadBundle(pb)
			}
		}
		changed, err := renderFile(filePath, dir, registry, engine, opts)
		processed++

		if err != nil {
			errorsList = append(errorsList, fmt.Errorf("%s: %w", filePath, err))
			continue
		}

		if changed {
			modified++
		}

		if opts.Verbose {
			logFileStatus(filePath, dir, opts.OutDir, changed)
		}
	}

	return &Result{
		FilesProcessed: processed,
		FilesModified:  modified,
		Errors:         errorsList,
		RenderErrors:   renderErrors,
	}, nil
}

// transformParallel processes files using a pool of QJS engines.
func transformParallel(htmlFiles []string, dir string, registry *jsengine.Registry, opts Options, preloadBundles []string, workers int) (*Result, error) {
	pool, err := jsengine.NewEnginePool(workers)
	if err != nil {
		return nil, fmt.Errorf("creating engine pool: %w", err)
	}
	defer pool.Close()

	// Load preload bundles into every engine in the pool.
	// We need to do this before PreloadAll since preload bundles
	// may register modules that component bundles depend on.
	drained := make([]*jsengine.Engine, workers)
	for i := 0; i < workers; i++ {
		e := pool.Get()
		e.SetPreloadModules(opts.Preload)
		for _, pb := range preloadBundles {
			_ = e.LoadBundle(pb)
		}
		drained[i] = e
	}
	for _, e := range drained {
		pool.Put(e)
	}

	// Pre-load all discovered component bundles into every engine.
	if err := pool.PreloadAll(registry, opts.Preload); err != nil {
		return nil, fmt.Errorf("preloading pool: %w", err)
	}

	type fileResult struct {
		filePath     string
		changed      bool
		err          error
		renderErrors []RenderError
	}

	results := make([]fileResult, len(htmlFiles))
	var wg sync.WaitGroup

	work := make(chan int, len(htmlFiles))
	for i := range htmlFiles {
		work <- i
	}
	close(work)

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			engine := pool.Get()
			defer pool.Put(engine)

			for i := range work {
				filePath := htmlFiles[i]
				changed, err := renderFile(filePath, dir, registry, engine, opts)
				results[i] = fileResult{
					filePath: filePath,
					changed:  changed,
					err:      err,
				}
			}
		}()
	}

	wg.Wait()

	var (
		processed    int
		modified     int
		errorsList   []error
		renderErrors []RenderError
	)

	for _, r := range results {
		processed++
		if r.err != nil {
			errorsList = append(errorsList, fmt.Errorf("%s: %w", r.filePath, r.err))
			continue
		}
		if r.changed {
			modified++
		}
		renderErrors = append(renderErrors, r.renderErrors...)
		if opts.Verbose {
			logFileStatus(r.filePath, dir, opts.OutDir, r.changed)
		}
	}

	return &Result{
		FilesProcessed: processed,
		FilesModified:  modified,
		Errors:         errorsList,
		RenderErrors:   renderErrors,
	}, nil
}

// logFileStatus prints a verbose status line for a processed file.
func logFileStatus(filePath, dir, outDir string, changed bool) {
	status := "unchanged"
	if changed {
		status = "modified"
	}
	outPath := filePath
	if outDir != "" {
		rel, _ := filepath.Rel(dir, filePath)
		outPath = filepath.Join(outDir, rel)
	}
	fmt.Fprintf(os.Stderr, "  %s -> %s [%s]\n", filePath, outPath, status)
}

// renderFile reads, transforms, and writes a single HTML file.
// Discovery must have already been run (pass 1); this function only renders.
func renderFile(filePath string, srcDir string, registry *jsengine.Registry, engine *jsengine.Engine, opts Options) (bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("reading file: %w", err)
	}

	input := string(data)

	ctx := &transformContext{
		engine:   engine,
		registry: registry,
		ignored:  opts.Ignored,
		file:     filePath,
	}
	output, err := renderHTMLWithContext(input, ctx)
	if err != nil {
		return false, fmt.Errorf("rendering: %w", err)
	}

	if output == input && opts.OutDir == "" {
		return false, nil
	}

	if !opts.DryRun {
		destPath := filePath
		if opts.OutDir != "" {
			rel, err := filepath.Rel(srcDir, filePath)
			if err != nil {
				return false, fmt.Errorf("computing relative path: %w", err)
			}
			destPath = filepath.Join(opts.OutDir, rel)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return false, fmt.Errorf("creating output directory: %w", err)
			}
		}
		if err := fileutil.WriteFileAtomic(destPath, []byte(output), 0644); err != nil {
			return false, fmt.Errorf("writing file: %w", err)
		}
	}

	changed := output != input
	if opts.OutDir != "" {
		changed = true
	}
	return changed, nil
}

// importRe matches ES module import statements to extract bare-module specifiers.
// Handles: import 'x', import "x", import {...} from 'x', import x from 'x'
var importRe = regexp.MustCompile(`import\s+(?:[^'"]*\s+from\s+)?['"]([^'"]+)['"]`)

// discoverFromHTML extracts import maps and module import specifiers from
// HTML content, resolves them, bundles the components, and registers them.
// siteRoot is the top-level directory passed to TransformDir (e.g. "public/"),
// used to resolve absolute paths like "/node_modules/..." in import maps.
func discoverFromHTML(htmlContent string, htmlDir string, siteRoot string, registry *jsengine.Registry, cliImportMap *jsengine.ImportMap, verbose bool) {
	// Parse the HTML to find script tags
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return
	}

	var htmlImportMap *jsengine.ImportMap
	var moduleSpecifiers []string

	// Walk the DOM to find script tags
	var walkNode func(*html.Node)
	walkNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			scriptType := getAttr(n, "type")

			if scriptType == "importmap" {
				// Extract the import map JSON from the script content.
				// Use siteRoot as baseDir so that absolute paths like
				// "/node_modules/..." resolve relative to the site root.
				content := getTextContent(n)
				if content != "" {
					absSiteRoot, _ := filepath.Abs(siteRoot)
					im, err := jsengine.ParseImportMap(content, absSiteRoot)
					if err == nil {
						htmlImportMap = im
					}
				}
			} else if scriptType == "module" {
				// Extract import specifiers from inline module scripts
				content := getTextContent(n)
				if content != "" {
					matches := importRe.FindAllStringSubmatch(content, -1)
					for _, match := range matches {
						if len(match) >= 2 {
							moduleSpecifiers = append(moduleSpecifiers, match[1])
						}
					}
				}
				// Also check src attribute for external module scripts
				// (we can't read those, but the inline ones are most common for import maps)
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walkNode(child)
		}
	}
	walkNode(doc)

	// CLI import map takes precedence, fall back to HTML-embedded map
	activeMap := cliImportMap
	if activeMap == nil {
		activeMap = htmlImportMap
	}

	if activeMap == nil || len(moduleSpecifiers) == 0 {
		if verbose {
			fmt.Fprintf(os.Stderr, "  golit: discovery: activeMap=%v specifiers=%v\n", activeMap != nil, moduleSpecifiers)
		}
		return
	}

	// Resolve specifiers through the import map and bundle
	resolvedPaths := activeMap.ResolveAll(moduleSpecifiers)
	if verbose {
		fmt.Fprintf(os.Stderr, "  golit: discovery: %d specifiers -> %d resolved paths\n", len(moduleSpecifiers), len(resolvedPaths))
		for _, p := range resolvedPaths {
			fmt.Fprintf(os.Stderr, "    %s\n", p)
		}
	}
	// Collect local paths for batch bundling; warn about CDN URLs
	var localPaths []string
	for _, path := range resolvedPaths {
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			fmt.Fprintf(os.Stderr, "  golit: skipping %s (CDN URL)\n", path)
			fmt.Fprintf(os.Stderr, "         Use --importmap with local paths or --sources for SSR\n")
			continue
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		localPaths = append(localPaths, path)
	}

	if len(localPaths) == 0 {
		return
	}

	// Skip bundling if all paths have already been processed.
	// This avoids re-running esbuild + QJS discovery for every HTML file
	// when they all share the same import map and module scripts.
	allKnown := true
	for _, p := range localPaths {
		if !registry.HasPath(p) {
			allKnown = false
			break
		}
	}
	if allKnown {
		return
	}

	// Batch-bundle all discovered components in one esbuild call
	bundles, err := jsengine.BundleComponents(localPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  golit: warning: batch bundle failed: %v\n", err)
		return
	}

	for path, bundle := range bundles {
		tagName, err := jsengine.DiscoverTagName(bundle)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  golit: warning: could not discover tag in %s: %v\n", path, err)
			registry.MarkPath(path) // mark as processed even if no tag found
			continue
		}

		if !registry.Has(tagName) {
			registry.Register(tagName, bundle)
			fmt.Fprintf(os.Stderr, "  golit: auto-discovered <%s> from %s\n", tagName, path)
		}
		registry.MarkPath(path)
	}
}

// getAttr gets an attribute value from an HTML node.
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// getTextContent gets the text content of an HTML node.
func getTextContent(n *html.Node) string {
	var buf strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			buf.WriteString(child.Data)
		}
	}
	return buf.String()
}

// RenderHTML takes an HTML string, finds custom elements, and returns
// the transformed HTML with Declarative Shadow DOM.
func RenderHTML(input string, registry *jsengine.Registry, ignored ...map[string]bool) (string, error) {
	var ign map[string]bool
	if len(ignored) > 0 {
		ign = ignored[0]
	}
	return renderHTMLWithIgnored(input, registry, ign)
}

func renderHTMLWithIgnored(input string, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	engine, err := jsengine.NewEngine()
	if err != nil {
		return "", fmt.Errorf("creating JS engine: %w", err)
	}
	defer engine.Close()
	return renderHTMLWithEngine(input, engine, registry, ignored)
}

// transformContext carries shared state through the recursive transform walk.
type transformContext struct {
	engine       *jsengine.Engine
	registry     *jsengine.Registry
	ignored      map[string]bool
	file         string // current HTML file path (for error reporting)
	renderErrors []RenderError
}

// renderHTMLWithEngine transforms HTML using a provided engine (no create/destroy overhead).
func renderHTMLWithEngine(input string, engine *jsengine.Engine, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	ctx := &transformContext{engine: engine, registry: registry, ignored: ignored}
	output, err := renderHTMLWithContext(input, ctx)
	return output, err
}

func renderHTMLWithContext(input string, ctx *transformContext) (string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", fmt.Errorf("parsing HTML: %w", err)
	}

	if err := renderHTMLBatched(doc, ctx, 10); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("rendering HTML: %w", err)
	}

	result := buf.String()
	if !isFullDocument(input) {
		result = extractBodyContent(result)
	}
	return result, nil
}

// RenderFragment renders an HTML fragment.
func RenderFragment(input string, registry *jsengine.Registry, ignored ...map[string]bool) (string, error) {
	var ign map[string]bool
	if len(ignored) > 0 {
		ign = ignored[0]
	}
	return renderFragmentWithIgnored(input, registry, ign)
}

func renderFragmentWithIgnored(input string, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	nodes, err := html.ParseFragment(strings.NewReader(input), &html.Node{
		Type: html.ElementNode, Data: "body", DataAtom: atom.Body,
	})
	if err != nil {
		return "", fmt.Errorf("parsing fragment: %w", err)
	}

	engine, err := jsengine.NewEngine()
	if err != nil {
		return "", fmt.Errorf("creating JS engine: %w", err)
	}
	defer engine.Close()

	ctx := &transformContext{engine: engine, registry: registry, ignored: ignored}

	// Wrap nodes in a temporary parent for batch rendering
	wrapper := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	for _, node := range nodes {
		wrapper.AppendChild(node)
	}

	if err := renderHTMLBatched(wrapper, ctx, 10); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	for child := wrapper.FirstChild; child != nil; child = child.NextSibling {
		if err := html.Render(&buf, child); err != nil {
			return "", fmt.Errorf("rendering: %w", err)
		}
	}
	return buf.String(), nil
}

// transformNode recursively walks the HTML tree and expands custom elements.
func transformNode(node *html.Node, ctx *transformContext, depth, maxDepth int) error {
	if depth > maxDepth {
		return nil
	}

	var originalChildren []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		originalChildren = append(originalChildren, child)
	}

	if node.Type == html.ElementNode && strings.Contains(node.Data, "-") {
		if ctx.ignored[node.Data] {
			// Explicitly ignored, skip SSR
		} else if hasDeclarativeShadowRoot(node) {
			// Already SSR'd, skip
		} else if ctx.engine.LoadBundleForTag(node.Data, ctx.registry) {
			if err := expandCustomElement(node, ctx, depth, maxDepth); err != nil {
				return fmt.Errorf("expanding <%s>: %w", node.Data, err)
			}
		} else if ctx.engine.IsRegistered(node.Data) {
			if err := expandCustomElement(node, ctx, depth, maxDepth); err != nil {
				// Sub-component render failed — already tracked in ctx.renderErrors
			}
		} else {
			ctx.registry.MarkUnregistered(node.Data)
		}
	}

	for _, child := range originalChildren {
		if err := transformNode(child, ctx, depth, maxDepth); err != nil {
			return err
		}
	}

	return nil
}

// expandCustomElement renders a component and adds DSD.
func expandCustomElement(node *html.Node, ctx *transformContext, depth, maxDepth int) error {
	attrs := make(map[string]string)
	for _, attr := range node.Attr {
		attrs[attr.Key] = attr.Val
	}

	result, err := ctx.engine.RenderElement(node.Data, attrs)
	if err != nil {
		ctx.renderErrors = append(ctx.renderErrors, RenderError{
			TagName: node.Data,
			File:    ctx.file,
			Err:     err,
		})
		return nil
	}

	// Build shadow DOM content: <style>CSS</style> + rendered HTML
	var shadowContent strings.Builder
	if result.CSS != "" {
		shadowContent.WriteString("<style>")
		shadowContent.WriteString(result.CSS)
		shadowContent.WriteString("</style>")
	}
	shadowContent.WriteString(result.HTML)

	// Parse the shadow content
	shadowNodes, err := html.ParseFragment(
		strings.NewReader(shadowContent.String()),
		&html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body},
	)
	if err != nil {
		return fmt.Errorf("parsing shadow content: %w", err)
	}

	// Create the <template> element
	templateNode := &html.Node{
		Type: html.ElementNode, Data: "template", DataAtom: atom.Template,
		Attr: []html.Attribute{
			{Key: "shadowroot", Val: "open"},
			{Key: "shadowrootmode", Val: "open"},
		},
	}

	for _, sn := range shadowNodes {
		templateNode.AppendChild(sn)
	}

	for child := templateNode.FirstChild; child != nil; child = child.NextSibling {
		if err := transformNode(child, ctx, depth+1, maxDepth); err != nil {
			return err
		}
	}

	// Add defer-hydration for nested custom elements
	if depth > 0 {
		node.Attr = append(node.Attr, html.Attribute{Key: "defer-hydration"})
	}

	// Insert the template as the first child
	if node.FirstChild != nil {
		node.InsertBefore(templateNode, node.FirstChild)
	} else {
		node.AppendChild(templateNode)
	}

	return nil
}

// pendingElement tracks a custom element waiting to be expanded in batch mode.
type pendingElement struct {
	node  *html.Node
	depth int
}

// collectUnexpanded walks the HTML tree and returns all custom elements
// that haven't been expanded yet (no <template shadowrootmode> child).
func collectUnexpanded(node *html.Node, ctx *transformContext) []pendingElement {
	var pending []pendingElement
	var walk func(*html.Node, int)
	walk = func(n *html.Node, depth int) {
		if depth > 10 {
			return
		}
		if n.Type == html.ElementNode && strings.Contains(n.Data, "-") {
			if !ctx.ignored[n.Data] && !hasDeclarativeShadowRoot(n) {
				if ctx.engine.LoadBundleForTag(n.Data, ctx.registry) || ctx.engine.IsRegistered(n.Data) {
					pending = append(pending, pendingElement{node: n, depth: depth})
					return // don't recurse into this node's children yet
				}
				ctx.registry.MarkUnregistered(n.Data)
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child, depth)
		}
	}
	walk(node, 0)
	return pending
}

// renderHTMLBatched uses BFS-by-depth to render all custom elements,
// batching all elements at each depth level into a single QJS Eval call.
func renderHTMLBatched(doc *html.Node, ctx *transformContext, maxDepth int) error {
	for depth := 0; depth < maxDepth; depth++ {
		pending := collectUnexpanded(doc, ctx)
		if len(pending) == 0 {
			break
		}

		requests := make([]jsengine.BatchRequest, len(pending))
		for i, p := range pending {
			attrs := make(map[string]string)
			for _, attr := range p.node.Attr {
				attrs[attr.Key] = attr.Val
			}
			requests[i] = jsengine.BatchRequest{
				ID:      i,
				TagName: p.node.Data,
				Attrs:   attrs,
			}
		}

		results, err := ctx.engine.RenderBatch(requests)
		if err != nil {
			return fmt.Errorf("batch render at depth %d: %w", depth, err)
		}

		resultMap := make(map[int]jsengine.BatchResult, len(results))
		for _, r := range results {
			resultMap[r.ID] = r
		}

		for i, p := range pending {
			r, ok := resultMap[i]
			if !ok {
				continue
			}
			if r.Error != "" {
				ctx.renderErrors = append(ctx.renderErrors, RenderError{
					TagName: p.node.Data,
					File:    ctx.file,
					Err:     fmt.Errorf("%s", r.Error),
				})
				continue
			}

			var shadowContent strings.Builder
			if r.CSS != "" {
				shadowContent.WriteString("<style>")
				shadowContent.WriteString(strings.TrimSpace(r.CSS))
				shadowContent.WriteString("</style>")
			}
			shadowContent.WriteString(r.HTML)

			shadowNodes, err := html.ParseFragment(
				strings.NewReader(shadowContent.String()),
				&html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body},
			)
			if err != nil {
				continue
			}

			templateNode := &html.Node{
				Type: html.ElementNode, Data: "template", DataAtom: atom.Template,
				Attr: []html.Attribute{
					{Key: "shadowroot", Val: "open"},
					{Key: "shadowrootmode", Val: "open"},
				},
			}

			for _, sn := range shadowNodes {
				templateNode.AppendChild(sn)
			}

			if depth > 0 {
				p.node.Attr = append(p.node.Attr, html.Attribute{Key: "defer-hydration"})
			}

			if p.node.FirstChild != nil {
				p.node.InsertBefore(templateNode, p.node.FirstChild)
			} else {
				p.node.AppendChild(templateNode)
			}
		}
	}
	return nil
}

// hasDeclarativeShadowRoot checks if an element already has DSD.
func hasDeclarativeShadowRoot(node *html.Node) bool {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "template" {
			for _, attr := range child.Attr {
				if attr.Key == "shadowrootmode" || attr.Key == "shadowroot" {
					return true
				}
			}
		}
	}
	return false
}

func isFullDocument(input string) bool {
	lower := strings.TrimSpace(strings.ToLower(input))
	return strings.HasPrefix(lower, "<!doctype") || strings.HasPrefix(lower, "<html")
}

func extractBodyContent(rendered string) string {
	bodyStart := strings.Index(rendered, "<body>")
	if bodyStart == -1 {
		return rendered
	}
	bodyStart += len("<body>")
	bodyEnd := strings.LastIndex(rendered, "</body>")
	if bodyEnd == -1 {
		return rendered
	}
	return rendered[bodyStart:bodyEnd]
}

func collectHTMLFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".html") ||
			strings.HasSuffix(strings.ToLower(info.Name()), ".htm") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
