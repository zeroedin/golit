// Package transformer provides file-level and directory-level HTML
// transformation, walking HTML files on disk and expanding custom
// elements using the QJS rendering engine.
package transformer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

	// Registry is an optional pre-populated bundle registry. When set,
	// TransformDir uses it instead of creating a new one. Additional
	// discovery modes (DefsDir, SourcesDir, etc.) still layer onto it.
	Registry *jsengine.Registry
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
	registry := opts.Registry
	if registry == nil {
		registry = jsengine.NewRegistry()
	}

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
	if !autoDiscover && opts.DefsDir == "" && opts.CompiledFile == "" && opts.SourcesDir == "" && opts.ImportMapFile == "" && opts.Registry == nil {
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
	if autoDiscover || cliImportMap != nil {
		for _, filePath := range htmlFiles {
			data, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			htmlDir := filepath.Dir(filePath)
			discoverFromHTML(string(data), htmlDir, dir, registry, cliImportMap, opts.Verbose)
		}
	}

	// ── Pass 2: Render ──────────────────────────────────────────────
	workers := opts.Concurrency
	if workers < 1 {
		workers = 1
	}
	if opts.Isolate {
		workers = 1
	}
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
	var checkEngine *jsengine.Engine
	if workers == 1 {
		result, checkEngine, err = transformSequential(htmlFiles, dir, registry, opts, preloadBundles)
	} else {
		result, checkEngine, err = transformParallel(htmlFiles, dir, registry, opts, preloadBundles, workers)
	}
	if err != nil {
		return nil, err
	}
	if checkEngine != nil {
		defer checkEngine.Close()
	}
	knownTags := registry.TagNames()
	knownPaths := registry.ProcessedPaths()
	var trulyUnregistered []string
	for _, tag := range registry.Unregistered() {
		if checkEngine != nil && checkEngine.IsRegistered(tag) {
			continue
		}

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

func transformSequential(htmlFiles []string, dir string, registry *jsengine.Registry, opts Options, preloadBundles []string) (*Result, *jsengine.Engine, error) {
	engine, err := initEngine(preloadBundles, opts.Preload)
	if err != nil {
		return nil, nil, fmt.Errorf("creating JS engine: %w", err)
	}

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
			engine.SetPreloadModules(opts.Preload)
			for _, pb := range preloadBundles {
				_ = engine.LoadBundle(pb)
			}
		}
		changed, reErrs, err := renderFile(filePath, dir, registry, engine, opts)
		processed++
		renderErrors = append(renderErrors, reErrs...)

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
	}, engine, nil
}

func transformParallel(htmlFiles []string, dir string, registry *jsengine.Registry, opts Options, preloadBundles []string, workers int) (*Result, *jsengine.Engine, error) {
	pool, err := jsengine.NewEnginePool(workers)
	if err != nil {
		return nil, nil, fmt.Errorf("creating engine pool: %w", err)
	}
	defer pool.Close()

	if err := pool.PreloadAll(registry, opts.Preload, preloadBundles...); err != nil {
		return nil, nil, fmt.Errorf("preloading pool: %w", err)
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
				changed, reErrs, err := renderFile(filePath, dir, registry, engine, opts)
				results[i] = fileResult{
					filePath:     filePath,
					changed:      changed,
					err:          err,
					renderErrors: reErrs,
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

	checkEngine := pool.Get()

	return &Result{
		FilesProcessed: processed,
		FilesModified:  modified,
		Errors:         errorsList,
		RenderErrors:   renderErrors,
	}, checkEngine, nil
}

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

func renderFile(filePath string, srcDir string, registry *jsengine.Registry, engine ElementRenderer, opts Options) (bool, []RenderError, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, nil, fmt.Errorf("reading file: %w", err)
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
		return false, ctx.renderErrors, fmt.Errorf("rendering: %w", err)
	}

	if output == input && opts.OutDir == "" {
		return false, ctx.renderErrors, nil
	}

	if !opts.DryRun {
		destPath := filePath
		if opts.OutDir != "" {
			rel, err := filepath.Rel(srcDir, filePath)
			if err != nil {
				return false, ctx.renderErrors, fmt.Errorf("computing relative path: %w", err)
			}
			destPath = filepath.Join(opts.OutDir, rel)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return false, ctx.renderErrors, fmt.Errorf("creating output directory: %w", err)
			}
		}
		if err := fileutil.WriteFileAtomic(destPath, []byte(output), 0644); err != nil {
			return false, ctx.renderErrors, fmt.Errorf("writing file: %w", err)
		}
	}

	changed := output != input
	if opts.OutDir != "" {
		changed = true
	}
	return changed, ctx.renderErrors, nil
}

func collectHTMLFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := strings.ToLower(d.Name())
		if strings.HasSuffix(name, ".html") || strings.HasSuffix(name, ".htm") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
