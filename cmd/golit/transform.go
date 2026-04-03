package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sspriggs/golit/pkg/transformer"
)

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
		case "--concurrency", "-j":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				n, err := strconv.Atoi(args[i+1])
				if err != nil {
					return fmt.Errorf("--concurrency value must be a positive integer")
				}
				if n < 1 {
					return fmt.Errorf("--concurrency value must be a positive integer")
				}
				cliOpts.Concurrency = n
				i += 2
			} else {
				cliOpts.Concurrency = runtime.NumCPU()
				i++
			}
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

	var opts transformer.Options
	if cfg != nil {
		opts = cfg.ToTransformOptions(cliOpts)
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
