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

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/sspriggs/golit/pkg/jsengine"
)

// Options configures the transformer.
type Options struct {
	// DefsDir is the directory containing .golit.bundle.js files (Mode 1).
	DefsDir string

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

	// Concurrency is the number of files to process in parallel.
	// 0 means use a sensible default.
	Concurrency int

	// OutDir is an optional output directory. When set, transformed files
	// are written here instead of modifying the input files in-place.
	OutDir string
}

// Result holds stats from a transform run.
type Result struct {
	FilesProcessed int
	FilesModified  int
	Errors         []error
	Unregistered   []string
}

// TransformDir processes all HTML files in a directory tree.
func TransformDir(dir string, opts Options) (*Result, error) {
	registry := jsengine.NewRegistry()

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
	if !autoDiscover && opts.DefsDir == "" && opts.SourcesDir == "" && opts.ImportMapFile == "" {
		autoDiscover = true
	}

	// Collect HTML files
	htmlFiles, err := collectHTMLFiles(dir)
	if err != nil {
		return nil, fmt.Errorf("collecting HTML files: %w", err)
	}

	// Create a single QJS engine for the entire transform run.
	// Bundles are loaded lazily via LoadBundleForTag and persist across files,
	// avoiding repeated engine creation and bundle re-evaluation per file.
	// QJS is not goroutine-safe, so files are processed sequentially.
	engine, err := jsengine.NewEngine()
	if err != nil {
		return nil, fmt.Errorf("creating JS engine: %w", err)
	}
	defer engine.Close()

	// Pre-load extra modules into QJS before component rendering.
	// These are resolved, bundled, and loaded so that dynamic import()
	// calls in component bundles can access them.
	// Tell the engine which modules are preloaded so it can shim import() calls.
	engine.SetPreloadModules(opts.Preload)

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
		if err := engine.LoadBundle(bundle); err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: could not load preload %s: %v\n", mod, err)
			continue
		}
		if opts.Verbose {
			fmt.Fprintf(os.Stderr, "golit: preloaded %s from %s\n", mod, modPath)
		}
	}

	var (
		processed  int
		modified   int
		errorsList []error
	)

	for _, filePath := range htmlFiles {
		changed, err := transformFile(filePath, dir, registry, engine, opts, cliImportMap, autoDiscover)
		processed++

		if err != nil {
			errorsList = append(errorsList, fmt.Errorf("%s: %w", filePath, err))
			continue
		}

		if changed {
			modified++
		}

		if opts.Verbose {
			status := "unchanged"
			if changed {
				status = "modified"
			}
			outPath := filePath
			if opts.OutDir != "" {
				rel, _ := filepath.Rel(dir, filePath)
				outPath = filepath.Join(opts.OutDir, rel)
			}
			fmt.Fprintf(os.Stderr, "  %s -> %s [%s]\n", filePath, outPath, status)
		}
	}

	// Filter the unregistered list: remove sub-components of known elements
	// and elements defined in the QJS engine. Sub-components are identified by:
	// 1. Prefix match: rh-accordion-header starts with rh-accordion-
	// 2. Source file match: rh-footer-block.js exists in rh-footer/ directory
	// 3. QJS registration: defined by a parent bundle's side effects
	knownTags := registry.TagNames()
	knownPaths := registry.ProcessedPaths()
	var trulyUnregistered []string
	for _, tag := range registry.Unregistered() {
		if engine.IsRegistered(tag) {
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
		// e.g. rh-footer-block.js in the same dir as rh-footer.js
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

	return &Result{
		FilesProcessed: processed,
		FilesModified:  modified,
		Errors:         errorsList,
		Unregistered:   trulyUnregistered,
	}, nil
}

// transformFile reads, transforms, and writes a single HTML file.
func transformFile(filePath string, srcDir string, registry *jsengine.Registry, engine *jsengine.Engine, opts Options, cliImportMap *jsengine.ImportMap, autoDiscover bool) (bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("reading file: %w", err)
	}

	input := string(data)

	// Mode 3 + 4: Discover components from import maps and module scripts
	htmlDir := filepath.Dir(filePath)
	if autoDiscover || cliImportMap != nil {
		discoverFromHTML(input, htmlDir, srcDir, registry, cliImportMap, opts.Verbose)
	}

	output, err := renderHTMLWithEngine(input, engine, registry, opts.Ignored)
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
		if err := os.WriteFile(destPath, []byte(output), 0644); err != nil {
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
func RenderHTML(input string, registry *jsengine.Registry) (string, error) {
	return renderHTMLWithIgnored(input, registry, nil)
}

func renderHTMLWithIgnored(input string, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	engine, err := jsengine.NewEngine()
	if err != nil {
		return "", fmt.Errorf("creating JS engine: %w", err)
	}
	defer engine.Close()
	return renderHTMLWithEngine(input, engine, registry, ignored)
}

// renderHTMLWithEngine transforms HTML using a provided engine (no create/destroy overhead).
func renderHTMLWithEngine(input string, engine *jsengine.Engine, registry *jsengine.Registry, ignored map[string]bool) (string, error) {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return "", fmt.Errorf("parsing HTML: %w", err)
	}

	if err := transformNode(doc, engine, registry, ignored, 0, 10); err != nil {
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
func RenderFragment(input string, registry *jsengine.Registry) (string, error) {
	return renderFragmentWithIgnored(input, registry, nil)
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

	for _, node := range nodes {
		if err := transformNode(node, engine, registry, ignored, 0, 10); err != nil {
			return "", err
		}
	}

	var buf bytes.Buffer
	for _, node := range nodes {
		if err := html.Render(&buf, node); err != nil {
			return "", fmt.Errorf("rendering: %w", err)
		}
	}
	return buf.String(), nil
}

// transformNode recursively walks the HTML tree and expands custom elements.
func transformNode(node *html.Node, engine *jsengine.Engine, registry *jsengine.Registry, ignored map[string]bool, depth, maxDepth int) error {
	if depth > maxDepth {
		return nil
	}

	var originalChildren []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		originalChildren = append(originalChildren, child)
	}

	if node.Type == html.ElementNode && strings.Contains(node.Data, "-") {
		if ignored[node.Data] {
			// Explicitly ignored, skip SSR
		} else if hasDeclarativeShadowRoot(node) {
			// Already SSR'd, skip
		} else if engine.LoadBundleForTag(node.Data, registry) {
			if err := expandCustomElement(node, engine, registry, ignored, depth, maxDepth); err != nil {
				return fmt.Errorf("expanding <%s>: %w", node.Data, err)
			}
		} else if engine.IsRegistered(node.Data) {
			// Sub-component registered by a parent bundle (e.g. rh-accordion-header
			// from rh-accordion). Try rendering it directly.
			if err := expandCustomElement(node, engine, registry, ignored, depth, maxDepth); err != nil {
				// Render failed (no render method, etc.) — that's fine, skip silently
			}
		} else {
			// Truly unknown custom element — not in registry and not defined in QJS
			registry.MarkUnregistered(node.Data)
		}
	}

	for _, child := range originalChildren {
		if err := transformNode(child, engine, registry, ignored, depth, maxDepth); err != nil {
			return err
		}
	}

	return nil
}

// expandCustomElement renders a component and adds DSD.
func expandCustomElement(node *html.Node, engine *jsengine.Engine, registry *jsengine.Registry, ignored map[string]bool, depth, maxDepth int) error {
	// Collect attributes from the HTML element
	attrs := make(map[string]string)
	for _, attr := range node.Attr {
		attrs[attr.Key] = attr.Val
	}

	// Render the element using QJS
	result, err := engine.RenderElement(node.Data, attrs)
	if err != nil {
		// On render error, leave the element as-is for client-side rendering
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

	// Recursively transform nested custom elements in the shadow DOM
	for child := templateNode.FirstChild; child != nil; child = child.NextSibling {
		if err := transformNode(child, engine, registry, ignored, depth+1, maxDepth); err != nil {
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
