package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sspriggs/golit/pkg/jsengine"
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
		if strings.HasSuffix(name, ".d.ts") {
			return nil
		}
		if strings.HasSuffix(name, ".golit.bundle.js") {
			return nil
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
