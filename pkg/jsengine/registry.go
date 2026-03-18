package jsengine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fastschema/qjs"
)

// Registry manages loaded component bundles and tracks which tag names
// are available for rendering.
type Registry struct {
	// bundles maps tag names to their bundle JS content
	bundles map[string]string

	// unregistered tracks custom element tags found but not in the registry
	unregistered map[string]bool

	// processedPaths tracks source file paths that have already been bundled,
	// so discoverFromHTML can skip re-bundling across multiple HTML files.
	processedPaths map[string]bool
}

// NewRegistry creates an empty bundle registry.
func NewRegistry() *Registry {
	return &Registry{
		bundles:        make(map[string]string),
		unregistered:   make(map[string]bool),
		processedPaths: make(map[string]bool),
	}
}

// LoadDir loads all .golit.bundle.js files from a directory.
// Each bundle is loaded into a temporary QJS context to discover which
// tag name it registers, then stored by that tag name.
func (r *Registry) LoadDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading bundles directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".golit.bundle.js") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if err := r.LoadFile(path); err != nil {
			return err
		}
	}

	return nil
}

// LoadFile loads a single .golit.bundle.js file and discovers its tag name.
// Bundles that don't register any custom elements are silently skipped.
func (r *Registry) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading bundle %s: %w", path, err)
	}

	bundle := string(data)

	// Discover the tag name by loading the bundle in a temporary QJS context
	tagName, err := discoverTagName(bundle)
	if err != nil {
		// Skip bundles that don't define a custom element (e.g. utility modules)
		return nil
	}

	r.bundles[tagName] = bundle
	return nil
}

// Register adds a bundle by tag name directly (for programmatic use).
func (r *Registry) Register(tagName string, bundle string) {
	r.bundles[tagName] = bundle
}

// Lookup returns the bundle JS for a given tag name, or "" if not found.
func (r *Registry) Lookup(tagName string) string {
	return r.bundles[tagName]
}

// Has returns true if a bundle is registered for the given tag name.
func (r *Registry) Has(tagName string) bool {
	_, ok := r.bundles[tagName]
	return ok
}

// TagNames returns all registered tag names.
func (r *Registry) TagNames() []string {
	names := make([]string, 0, len(r.bundles))
	for name := range r.bundles {
		names = append(names, name)
	}
	return names
}

// MarkUnregistered records a custom element tag that was encountered
// but not found in the registry.
func (r *Registry) MarkUnregistered(tagName string) {
	r.unregistered[tagName] = true
}

// Unregistered returns all custom element tags that were encountered
// but not in the registry.
func (r *Registry) Unregistered() []string {
	tags := make([]string, 0, len(r.unregistered))
	for tag := range r.unregistered {
		tags = append(tags, tag)
	}
	return tags
}

// HasPath returns true if a source file path has already been processed.
func (r *Registry) HasPath(path string) bool {
	return r.processedPaths[path]
}

// MarkPath records a source file path as processed.
func (r *Registry) MarkPath(path string) {
	r.processedPaths[path] = true
}

// ProcessedPaths returns all source file paths that have been processed.
func (r *Registry) ProcessedPaths() []string {
	paths := make([]string, 0, len(r.processedPaths))
	for p := range r.processedPaths {
		paths = append(paths, p)
	}
	return paths
}

// LoadSourceDir bundles all .js/.ts files in a directory tree and registers them.
// This is Mode 2: the user points at a source directory.
// The directory is walked recursively so component files in subdirectories
// (e.g. node_modules/@rhds/elements/elements/rh-badge/rh-badge.js) are found.
// All files are bundled in a single esbuild invocation for performance.
func (r *Registry) LoadSourceDir(dir string) error {
	// Phase 1: collect all source file paths
	var paths []string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip inaccessible paths
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if strings.HasSuffix(name, ".d.ts") || strings.HasSuffix(name, ".golit.bundle.js") {
			return nil // skip declaration files and already-bundled files
		}
		ext := filepath.Ext(name)
		if ext != ".js" && ext != ".ts" && ext != ".tsx" {
			return nil
		}
		paths = append(paths, path)
		return nil
	}); err != nil {
		return fmt.Errorf("walking sources directory %s: %w", dir, err)
	}

	if len(paths) == 0 {
		return nil
	}

	// Phase 2: batch-bundle all files in one esbuild call
	bundles, err := BundleComponents(paths)
	if err != nil {
		return fmt.Errorf("batch bundling sources: %w", err)
	}

	// Phase 3: discover tag names and register
	for _, bundle := range bundles {
		tagName, err := DiscoverTagName(bundle)
		if err != nil {
			continue // skip files that don't define a custom element
		}
		r.bundles[tagName] = bundle
	}

	return nil
}

// LoadCompiled loads a single pre-compiled .golit.compiled.js artifact
// that contains all bundles and a __golitRegistry manifest.
func (r *Registry) LoadCompiled(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading compiled artifact %s: %w", path, err)
	}

	content := string(data)

	engine, err := NewEngine()
	if err != nil {
		return err
	}
	defer engine.Close()

	if err := engine.LoadBundle(content); err != nil {
		return fmt.Errorf("loading compiled artifact: %w", err)
	}

	result, err := engine.ctx.Eval("registry.js", qjs.Code(`
		(function() {
			if (globalThis.__golitRegistry) {
				return JSON.stringify(Object.keys(globalThis.__golitRegistry));
			}
			var reg = customElements;
			if (reg && reg.__definitions) {
				var names = [];
				for (var entry of reg.__definitions) { names.push(entry[0]); }
				return JSON.stringify(names);
			}
			return '[]';
		})();
	`))
	if err != nil {
		return fmt.Errorf("querying compiled registry: %w", err)
	}

	var names []string
	if err := json.Unmarshal([]byte(result.String()), &names); err != nil {
		return fmt.Errorf("parsing compiled registry: %w", err)
	}

	for _, name := range names {
		r.bundles[name] = content
	}

	return nil
}

// DiscoverTagName loads a bundle in a temporary QJS context and checks
// which tag name was registered via customElements.define().
func DiscoverTagName(bundle string) (string, error) {
	return discoverTagName(bundle)
}

func discoverTagName(bundle string) (string, error) {
	engine, err := NewEngine()
	if err != nil {
		return "", err
	}
	defer engine.Close()

	if err := engine.LoadBundle(bundle); err != nil {
		return "", err
	}

	// Query the custom elements registry for defined elements
	result, err := engine.ctx.Eval("discover.js", qjs.Code(`
		(function() {
			// Our DOM shim's customElements stores definitions in __definitions
			const reg = customElements;
			if (reg && reg.__definitions) {
				const names = [];
				for (const [name] of reg.__definitions) {
					names.push(name);
				}
				return JSON.stringify(names);
			}
			return '[]';
		})();
	`))
	if err != nil {
		return "", fmt.Errorf("querying tag names: %w", err)
	}

	var names []string
	if err := json.Unmarshal([]byte(result.String()), &names); err != nil {
		return "", fmt.Errorf("parsing tag names: %w", err)
	}

	if len(names) == 0 {
		return "", fmt.Errorf("no custom elements registered in bundle")
	}

	// Return the last registered name (typically the main component)
	return names[len(names)-1], nil
}
