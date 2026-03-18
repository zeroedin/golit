// golit is a CLI tool for server-side rendering Lit web components
// into Declarative Shadow DOM HTML using a QJS JavaScript engine.
//
// Usage:
//
//	golit bundle <source.ts|js> [--out <file>] [--minify]
//	golit transform <html-dir> --defs <bundles-dir> [--out <dir>] [--verbose]
//	golit render --defs <bundles-dir> '<html-fragment>'
//	golit version
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sspriggs/golit/pkg/jsengine"
	"github.com/sspriggs/golit/pkg/transformer"
)

const version = "0.2.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "bundle":
		if err := runBundle(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "compile":
		if err := runCompile(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "transform":
		if err := runTransform(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "render":
		if err := runRender(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version":
		fmt.Printf("golit %s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `golit - Lit SSR in pure Go (QJS engine)

Usage:
  golit bundle <source.ts|js> [--out <file.golit.bundle.js>] [--minify]
  golit bundle <src-dir/> [--out <bundles-dir/>] [--minify]
  golit compile --defs <bundles-dir> [--out <file.golit.compiled.js>] [--minify]
  golit transform <html-dir> [--defs <dir>] [--compiled <file>] [--sources <dir>] [--importmap <file>] [--out <dir>]
  golit render --defs <bundles-dir> '<html-fragment>'
  golit render --component-js '<source>' '<html-fragment>'
  golit version

Commands:
  bundle      Bundle a Lit component source file into a .golit.bundle.js
              file for SSR rendering. Includes Lit, DOM shim, and all deps.
  compile     Combine all bundles from a --defs directory into a single
              .golit.compiled.js artifact with a tag registry manifest.
  transform   Post-process HTML files, expanding custom elements into
              Declarative Shadow DOM using bundled components.
  render      Render a single HTML fragment to stdout.
  version     Print version information.

Options:
  --out <path>       Output path for bundle/transform
  --defs <dir>       Directory containing pre-bundled .golit.bundle.js files
  --compiled <file>  Single pre-compiled .golit.compiled.js artifact
  --sources <dir>    Directory of component .js/.ts source files (auto-bundles)
  --importmap <file> Import map JSON file for resolving bare-module specifiers
  --ignore <tag>     Skip SSR for this custom element (repeatable)
  --minify           Minify the output bundle
  --verbose          Print progress to stderr
  --dry-run          Process files without writing changes
  --strict           Exit with error if any components fail to render
  --component-js     Inline JS/TS component source for render command (repeatable)

Component Discovery (transform command):
  golit discovers which components to SSR using four modes (combinable):
  1. --defs      Pre-bundled .golit.bundle.js files
  2. --sources   Directory of source files (bundled on-demand)
  3. --importmap Import map file + module scripts in HTML
  4. Auto        Parse <script type="importmap"> and <script type="module">
                 from the HTML itself (default when no flags given)

Examples:
  # Auto-discover from HTML (zero config)
  hugo build && golit transform public/

  # Use pre-bundled components
  golit bundle src/components/ --out bundles/
  golit transform public/ --defs bundles/

  # Point at installed package sources
  golit transform public/ --sources node_modules/@rhds/elements/elements/

  # Use an import map
  golit transform public/ --importmap importmap.json

  # Render a single element
  golit render --defs bundles/ '<my-greeting name="World"></my-greeting>'
`)
}

// --- bundle command ---

func runBundle(args []string) error {
	var source, outPath string
	var opts jsengine.BundleOptions

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--out", "-o":
			if i+1 >= len(args) {
				return fmt.Errorf("--out requires a path argument")
			}
			outPath = args[i+1]
			i += 2
		case "--minify":
			opts.Minify = true
			i++
		default:
			if strings.HasPrefix(args[i], "--") {
				return fmt.Errorf("unknown option: %s", args[i])
			}
			if source == "" {
				source = args[i]
			} else {
				return fmt.Errorf("unexpected argument: %s", args[i])
			}
			i++
		}
	}

	if source == "" {
		return fmt.Errorf("missing required <source> argument")
	}

	info, err := os.Stat(source)
	if os.IsNotExist(err) {
		return fmt.Errorf("source does not exist: %s", source)
	}

	if info.IsDir() {
		return bundleDir(source, outPath, opts)
	}
	return bundleFile(source, outPath, opts)
}

func bundleFile(source, outPath string, opts jsengine.BundleOptions) error {
	bundle, err := jsengine.BundleComponent(source, opts)
	if err != nil {
		return err
	}

	if outPath == "" {
		ext := filepath.Ext(source)
		outPath = strings.TrimSuffix(source, ext) + ".golit.bundle.js"
	} else {
		// If outPath is a directory, generate the bundle filename inside it
		info, err := os.Stat(outPath)
		if err == nil && info.IsDir() {
			base := filepath.Base(source)
			ext := filepath.Ext(base)
			outPath = filepath.Join(outPath, strings.TrimSuffix(base, ext)+".golit.bundle.js")
		}
	}

	if err := jsengine.SaveBundle(bundle, outPath); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "golit: bundled %s -> %s (%d bytes)\n", source, outPath, len(bundle))
	return nil
}

func bundleDir(srcDir, outDir string, opts jsengine.BundleOptions) error {
	if outDir == "" {
		outDir = srcDir
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Collect all source files recursively, skipping declaration files
	var paths []string
	if err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip inaccessible paths
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if strings.HasSuffix(name, ".d.ts") {
			return nil // skip TypeScript declaration files
		}
		if strings.HasSuffix(name, ".golit.bundle.js") {
			return nil // skip already-bundled files
		}
		ext := filepath.Ext(name)
		if ext != ".js" && ext != ".ts" && ext != ".tsx" {
			return nil
		}
		paths = append(paths, path)
		return nil
	}); err != nil {
		return fmt.Errorf("walking source directory: %w", err)
	}

	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, "golit: 0 components bundled\n")
		return nil
	}

	// Batch-bundle all files in one esbuild call for performance
	bundles, err := jsengine.BundleComponents(paths, opts)
	if err != nil {
		return fmt.Errorf("batch bundling: %w", err)
	}

	count := 0
	for srcPath, bundle := range bundles {
		base := filepath.Base(srcPath)
		ext := filepath.Ext(base)
		outName := strings.TrimSuffix(base, ext) + ".golit.bundle.js"
		outPath := filepath.Join(outDir, outName)

		if err := jsengine.SaveBundle(bundle, outPath); err != nil {
			return fmt.Errorf("saving %s: %w", outPath, err)
		}

		fmt.Fprintf(os.Stderr, "golit: bundled %s -> %s (%d bytes)\n", srcPath, outPath, len(bundle))
		count++
	}

	fmt.Fprintf(os.Stderr, "golit: %d components bundled\n", count)
	return nil
}

// --- compile command ---

func runCompile(args []string) error {
	var defsDir, outPath string
	var minify bool

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--defs":
			if i+1 >= len(args) {
				return fmt.Errorf("--defs requires a directory argument")
			}
			defsDir = args[i+1]
			i += 2
		case "--out", "-o":
			if i+1 >= len(args) {
				return fmt.Errorf("--out requires a file argument")
			}
			outPath = args[i+1]
			i += 2
		case "--minify":
			minify = true
			i++
		default:
			if strings.HasPrefix(args[i], "--") {
				return fmt.Errorf("unknown option: %s", args[i])
			}
			if defsDir == "" {
				defsDir = args[i]
			}
			i++
		}
	}

	if defsDir == "" {
		return fmt.Errorf("missing required --defs <dir> argument")
	}
	if outPath == "" {
		outPath = "golit.compiled.js"
	}

	registry := jsengine.NewRegistry()
	if err := registry.LoadDir(defsDir); err != nil {
		return fmt.Errorf("loading bundles: %w", err)
	}

	tagNames := registry.TagNames()
	if len(tagNames) == 0 {
		return fmt.Errorf("no components found in %s", defsDir)
	}

	var compiled strings.Builder

	seen := make(map[string]bool)
	for _, tag := range tagNames {
		bundle := registry.Lookup(tag)
		if seen[bundle] {
			continue
		}
		seen[bundle] = true
		compiled.WriteString(bundle)
		compiled.WriteString("\n")
	}

	compiled.WriteString("globalThis.__golitRegistry = {")
	for i, tag := range tagNames {
		if i > 0 {
			compiled.WriteString(", ")
		}
		fmt.Fprintf(&compiled, "%q: true", tag)
	}
	compiled.WriteString("};\n")

	output := compiled.String()
	_ = minify // reserved for future esbuild minification pass

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(output), 0644); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "golit: compiled %d component(s) -> %s (%d bytes)\n", len(tagNames), outPath, len(output))
	return nil
}

// --- transform command ---

func runTransform(args []string) error {
	cliOpts := transformer.Options{}
	var htmlDir string
	var configPath string
	var strict bool

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--strict":
			strict = true
			i++
		case "--config", "-c":
			if i+1 >= len(args) {
				return fmt.Errorf("--config requires a file argument")
			}
			configPath = args[i+1]
			i += 2
		case "--defs":
			if i+1 >= len(args) {
				return fmt.Errorf("--defs requires a directory argument")
			}
			cliOpts.DefsDir = args[i+1]
			i += 2
		case "--compiled":
			if i+1 >= len(args) {
				return fmt.Errorf("--compiled requires a file argument")
			}
			cliOpts.CompiledFile = args[i+1]
			i += 2
		case "--sources":
			if i+1 >= len(args) {
				return fmt.Errorf("--sources requires a directory argument")
			}
			cliOpts.SourcesDir = args[i+1]
			i += 2
		case "--importmap":
			if i+1 >= len(args) {
				return fmt.Errorf("--importmap requires a file argument")
			}
			cliOpts.ImportMapFile = args[i+1]
			i += 2
		case "--ignore":
			if i+1 >= len(args) {
				return fmt.Errorf("--ignore requires a tag name argument")
			}
			if cliOpts.Ignored == nil {
				cliOpts.Ignored = make(map[string]bool)
			}
			cliOpts.Ignored[args[i+1]] = true
			i += 2
		case "--preload":
			if i+1 >= len(args) {
				return fmt.Errorf("--preload requires a module name argument")
			}
			cliOpts.Preload = append(cliOpts.Preload, args[i+1])
			i += 2
		case "--verbose", "-v":
			cliOpts.Verbose = true
			i++
		case "--dry-run":
			cliOpts.DryRun = true
			i++
		case "--isolate":
			cliOpts.Isolate = true
			i++
		case "--out", "-o":
			if i+1 >= len(args) {
				return fmt.Errorf("--out requires a directory argument")
			}
			cliOpts.OutDir = args[i+1]
			i += 2
		default:
			if strings.HasPrefix(args[i], "--") {
				return fmt.Errorf("unknown option: %s", args[i])
			}
			if htmlDir == "" {
				htmlDir = args[i]
			} else {
				return fmt.Errorf("unexpected argument: %s", args[i])
			}
			i++
		}
	}

	// Load config file: explicit --config, or auto-detect golit.yaml
	var cfg *Config
	if configPath != "" {
		var err error
		cfg, err = LoadConfig(configPath)
		if err != nil {
			return err
		}
	} else if found := FindConfig(); found != "" {
		var err error
		cfg, err = LoadConfig(found)
		if err != nil {
			return err
		}
	}

	// Merge config with CLI flags (CLI wins)
	var opts transformer.Options
	if cfg != nil {
		opts = cfg.ToTransformOptions(cliOpts)
		// Config can also provide the input directory
		if htmlDir == "" && cfg.Transform.Input != "" {
			htmlDir = cfg.Transform.Input
		}
	} else {
		opts = cliOpts
	}

	if htmlDir == "" {
		return fmt.Errorf("missing required <html-dir> argument")
	}

	start := time.Now()

	if opts.Verbose {
		if opts.OutDir != "" {
			fmt.Fprintf(os.Stderr, "golit transform: processing %s -> %s with bundles from %s\n", htmlDir, opts.OutDir, opts.DefsDir)
		} else {
			fmt.Fprintf(os.Stderr, "golit transform: processing %s (in-place) with bundles from %s\n", htmlDir, opts.DefsDir)
		}
	}

	result, err := transformer.TransformDir(htmlDir, opts)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "golit: %d files processed, %d modified in %s\n",
		result.FilesProcessed, result.FilesModified, elapsed.Round(time.Millisecond))

	if len(result.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "golit: %d errors:\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Fprintf(os.Stderr, "  - %v\n", err)
		}
	}

	if len(result.RenderErrors) > 0 {
		fmt.Fprintf(os.Stderr, "golit: %d component(s) failed to render (left as-is for client-side):\n", len(result.RenderErrors))
		for _, re := range result.RenderErrors {
			fmt.Fprintf(os.Stderr, "  - %s\n", re.Error())
		}
	}

	if len(result.Unregistered) > 0 {
		fmt.Fprintf(os.Stderr, "golit: %d custom element(s) found without bundles (passed through for client-side rendering):\n", len(result.Unregistered))
		for _, tag := range result.Unregistered {
			fmt.Fprintf(os.Stderr, "  - <%s>\n", tag)
		}
	}

	if strict && len(result.RenderErrors) > 0 {
		return fmt.Errorf("%d component(s) failed to render", len(result.RenderErrors))
	}

	return nil
}

// --- render command ---

func runRender(args []string) error {
	var defsDir, fragment string
	var componentSources []string

	i := 0
	for i < len(args) {
		switch args[i] {
		case "--defs":
			if i+1 >= len(args) {
				return fmt.Errorf("--defs requires a directory argument")
			}
			defsDir = args[i+1]
			i += 2
		case "--component-js":
			if i+1 >= len(args) {
				return fmt.Errorf("--component-js requires a JS source argument")
			}
			componentSources = append(componentSources, args[i+1])
			i += 2
		default:
			if strings.HasPrefix(args[i], "--") {
				return fmt.Errorf("unknown option: %s", args[i])
			}
			if fragment == "" {
				fragment = args[i]
			} else {
				fragment += " " + args[i]
			}
			i++
		}
	}

	if defsDir == "" && len(componentSources) == 0 {
		return fmt.Errorf("missing required --defs <dir> or --component-js <source> argument")
	}
	if fragment == "" {
		return fmt.Errorf("missing HTML fragment argument")
	}

	registry := jsengine.NewRegistry()
	if defsDir != "" {
		if err := registry.LoadDir(defsDir); err != nil {
			return fmt.Errorf("loading bundles: %w", err)
		}
	}

	for _, src := range componentSources {
		bundle, err := jsengine.BundleSource(src)
		if err != nil {
			return fmt.Errorf("bundling inline component: %w", err)
		}
		tagName, err := jsengine.DiscoverTagName(bundle)
		if err != nil {
			return fmt.Errorf("discovering tag from inline component: %w", err)
		}
		registry.Register(tagName, bundle)
		fmt.Fprintf(os.Stderr, "golit: registered <%s> from inline source\n", tagName)
	}

	output, err := transformer.RenderFragment(fragment, registry)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	fmt.Print(output)
	return nil
}
