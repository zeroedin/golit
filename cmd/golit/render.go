package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zeroedin/golit/pkg/jsengine"
	"github.com/zeroedin/golit/pkg/transformer"
)

func runRender(args []string) error {
	var defsDir, inputHTML string
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
			if inputHTML == "" {
				inputHTML = args[i]
			} else {
				inputHTML += " " + args[i]
			}
			i++
		}
	}

	if defsDir == "" && len(componentSources) == 0 {
		return fmt.Errorf("missing required --defs <dir> or --component-js <source> argument")
	}
	if inputHTML == "" {
		info, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("checking stdin: %w", err)
		}
		if info.Mode()&os.ModeCharDevice == 0 {
			const maxStdin = 32 << 20 // 32 MiB
			data, err := io.ReadAll(io.LimitReader(os.Stdin, maxStdin+1))
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			if len(data) > maxStdin {
				return fmt.Errorf("stdin input too large (max %d MiB)", maxStdin>>20)
			}
			inputHTML = strings.TrimSpace(string(data))
		}
	}
	if inputHTML == "" {
		return fmt.Errorf("missing HTML input (pass as argument or pipe to stdin)")
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

	output, err := transformer.RenderHTML(inputHTML, registry)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	fmt.Print(output)
	return nil
}
