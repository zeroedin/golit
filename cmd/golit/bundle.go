package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/zeroedin/golit/pkg/jsengine"
)

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
		case "--target":
			if i+1 >= len(args) {
				return fmt.Errorf("--target requires a value")
			}
			t, err := parseTarget(args[i+1])
			if err != nil {
				return err
			}
			opts.Target = t
			i += 2
		case "--format":
			if i+1 >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			f, err := parseFormat(args[i+1])
			if err != nil {
				return err
			}
			opts.Format = f
			i += 2
		case "--platform":
			if i+1 >= len(args) {
				return fmt.Errorf("--platform requires a value")
			}
			p, err := parsePlatform(args[i+1])
			if err != nil {
				return err
			}
			opts.Platform = p
			i += 2
		case "--conditions":
			if i+1 >= len(args) {
				return fmt.Errorf("--conditions requires a value")
			}
			opts.Conditions = strings.Split(args[i+1], ",")
			i += 2
		case "--main-fields":
			if i+1 >= len(args) {
				return fmt.Errorf("--main-fields requires a value")
			}
			opts.MainFields = strings.Split(args[i+1], ",")
			i += 2
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
		return bundleDirWithModules(source, outPath, opts)
	}
	return bundleFile(source, outPath, opts)
}

func bundleFile(source, outPath string, opts jsengine.BundleOptions) error {
	mod, err := jsengine.BundleComponentModule(source, opts)
	if err != nil {
		return err
	}

	if outPath == "" {
		ext := filepath.Ext(source)
		outPath = strings.TrimSuffix(source, ext) + ".golit.module.js"
	} else {
		info, err := os.Stat(outPath)
		if err == nil && info.IsDir() {
			base := filepath.Base(source)
			ext := filepath.Ext(base)
			outPath = filepath.Join(outPath, strings.TrimSuffix(base, ext)+".golit.module.js")
		}
	}

	if err := jsengine.SaveBundle(mod, outPath); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "golit: module %s -> %s (%d bytes)\n", source, outPath, len(mod))
	return nil
}

func bundleDirWithModules(srcDir, outDir string, opts jsengine.BundleOptions) error {
	if outDir == "" {
		outDir = srcDir
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Collect source files
	var paths []string
	if err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: skipping %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasSuffix(name, ".golit.bundle.js") || strings.HasSuffix(name, ".golit.module.js") {
			return nil
		}
		ext := filepath.Ext(name)
		if ext != ".js" {
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

	nodePaths := jsengine.FindAllNodeModules(paths[0])
	if len(nodePaths) == 0 {
		return fmt.Errorf("node_modules not found from %s", paths[0])
	}

	externals, err := jsengine.DiscoverExternalPackages(paths, nodePaths, opts)
	if err != nil {
		return fmt.Errorf("discovering external packages: %w", err)
	}
	opts.ExternalPackages = externals

	modules, err := jsengine.BundleComponentModules(paths, opts)
	if err != nil {
		return fmt.Errorf("batch bundling modules: %w", err)
	}

	runtime, err := jsengine.BundleSharedRuntime(nodePaths, modules, opts)
	if err != nil {
		return fmt.Errorf("building shared runtime: %w", err)
	}

	runtimePath := filepath.Join(outDir, "_runtime.golit.module.js")
	if err := jsengine.SaveBundle(runtime, runtimePath); err != nil {
		return fmt.Errorf("saving runtime: %w", err)
	}
	fmt.Fprintf(os.Stderr, "golit: shared runtime -> %s (%d bytes)\n", runtimePath, len(runtime))

	modules = jsengine.RewriteModuleImports(modules, externals)

	count := 0
	for srcPath, mod := range modules {
		base := filepath.Base(srcPath)
		ext := filepath.Ext(base)
		outName := strings.TrimSuffix(base, ext) + ".golit.module.js"
		outPath := filepath.Join(outDir, outName)

		if err := jsengine.SaveBundle(mod, outPath); err != nil {
			return fmt.Errorf("saving %s: %w", outPath, err)
		}

		fmt.Fprintf(os.Stderr, "golit: module %s -> %s (%d bytes)\n", srcPath, outPath, len(mod))
		count++
	}

	fmt.Fprintf(os.Stderr, "golit: %d component modules + 1 shared runtime bundled\n", count)
	return nil
}

func parseTarget(s string) (api.Target, error) {
	switch strings.ToLower(s) {
	case "esnext":
		return api.ESNext, nil
	case "es2015":
		return api.ES2015, nil
	case "es2016":
		return api.ES2016, nil
	case "es2017":
		return api.ES2017, nil
	case "es2018":
		return api.ES2018, nil
	case "es2019":
		return api.ES2019, nil
	case "es2020":
		return api.ES2020, nil
	case "es2021":
		return api.ES2021, nil
	case "es2022":
		return api.ES2022, nil
	case "es2023":
		return api.ES2023, nil
	case "es2024":
		return api.ES2024, nil
	default:
		return 0, fmt.Errorf("unknown target: %s (valid: esnext, es2015-es2024)", s)
	}
}

func parseFormat(s string) (api.Format, error) {
	switch strings.ToLower(s) {
	case "esm":
		return api.FormatESModule, nil
	case "cjs":
		return api.FormatCommonJS, nil
	case "iife":
		return api.FormatIIFE, nil
	default:
		return 0, fmt.Errorf("unknown format: %s (valid: esm, cjs, iife)", s)
	}
}

func parsePlatform(s string) (api.Platform, error) {
	switch strings.ToLower(s) {
	case "neutral":
		return api.PlatformNeutral, nil
	case "browser":
		return api.PlatformBrowser, nil
	case "node":
		return api.PlatformNode, nil
	default:
		return 0, fmt.Errorf("unknown platform: %s (valid: neutral, browser, node)", s)
	}
}
