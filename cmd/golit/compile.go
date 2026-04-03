package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zeroedin/golit/pkg/jsengine"
)

func runCompile(args []string) error {
	var defsDir, outPath string

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
	sort.Strings(tagNames)

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

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(output), 0644); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "golit: compiled %d component(s) -> %s (%d bytes)\n", len(tagNames), outPath, len(output))
	return nil
}
