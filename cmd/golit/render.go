package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zeroedin/golit/pkg/jsengine"
	"github.com/zeroedin/golit/pkg/transformer"
)

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
