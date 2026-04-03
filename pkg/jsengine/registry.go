package jsengine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fastschema/qjs"
)

// defineRe extracts custom element tag names from customElements.define() calls
// without executing the bundle. Valid custom element names must contain a hyphen.
var defineRe = regexp.MustCompile(`customElements\s*\.\s*define\s*\(\s*['"]([a-z][a-z0-9]*(?:-[a-z0-9]+)+)['"]`)

// discoveryJS is the snippet evaluated inside QJS to read registered tag names
// from the DOM shim's customElements.__definitions map.
const discoveryJS = `(function() {
	const reg = customElements;
	if (reg && reg.__definitions) {
		const names = [];
		for (const [name] of reg.__definitions) {
			names.push(name);
		}
		return JSON.stringify(names);
	}
	return '[]';
})();`

// Registry manages loaded component bundles and tracks which tag names
// are available for rendering. All methods are safe for concurrent use.
type Registry struct {
	mu sync.RWMutex

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
// Tag names are discovered via a regex pre-pass; bundles the regex
// misses are batched through a single reusable QJS engine.
func (r *Registry) LoadDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading bundles directory %s: %w", dir, err)
	}

	var bundles []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".golit.bundle.js") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("reading bundle %s: %w", entry.Name(), err)
		}
		bundles = append(bundles, string(data))
	}

	return r.registerBundles(bundles)
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

	r.mu.Lock()
	r.bundles[tagName] = bundle
	r.mu.Unlock()
	return nil
}

// registerBundles discovers tag names for a slice of bundles and registers
// them. Uses the regex fast path first, then batches remaining bundles
// through a single QJS engine.
func (r *Registry) registerBundles(bundles []string) error {
	var fallbackBundles []string

	r.mu.Lock()
	for _, bundle := range bundles {
		if tagName, ok := discoverTagNameFast(bundle); ok {
			r.bundles[tagName] = bundle
		} else {
			fallbackBundles = append(fallbackBundles, bundle)
		}
	}
	r.mu.Unlock()

	if len(fallbackBundles) > 0 {
		discovered, err := discoverTagNames(fallbackBundles)
		if err != nil {
			return err
		}
		r.mu.Lock()
		for j, tagName := range discovered {
			r.bundles[tagName] = fallbackBundles[j]
		}
		r.mu.Unlock()
	}

	return nil
}

// Register adds a bundle by tag name directly (for programmatic use).
func (r *Registry) Register(tagName string, bundle string) {
	r.mu.Lock()
	r.bundles[tagName] = bundle
	r.mu.Unlock()
}

// Lookup returns the bundle JS for a given tag name, or "" if not found.
func (r *Registry) Lookup(tagName string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.bundles[tagName]
}

// Has returns true if a bundle is registered for the given tag name.
func (r *Registry) Has(tagName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.bundles[tagName]
	return ok
}

// TagNames returns all registered tag names.
func (r *Registry) TagNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.bundles))
	for name := range r.bundles {
		names = append(names, name)
	}
	return names
}

// MarkUnregistered records a custom element tag that was encountered
// but not found in the registry.
func (r *Registry) MarkUnregistered(tagName string) {
	r.mu.Lock()
	r.unregistered[tagName] = true
	r.mu.Unlock()
}

// Unregistered returns all custom element tags that were encountered
// but not in the registry.
func (r *Registry) Unregistered() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tags := make([]string, 0, len(r.unregistered))
	for tag := range r.unregistered {
		tags = append(tags, tag)
	}
	return tags
}

// HasPath returns true if a source file path has already been processed.
func (r *Registry) HasPath(path string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.processedPaths[path]
}

// MarkPath records a source file path as processed.
func (r *Registry) MarkPath(path string) {
	r.mu.Lock()
	r.processedPaths[path] = true
	r.mu.Unlock()
}

// ProcessedPaths returns all source file paths that have been processed.
func (r *Registry) ProcessedPaths() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
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

	// Phase 3: discover tag names and register (regex fast path + batched QJS fallback)
	bundleList := make([]string, 0, len(bundles))
	for _, b := range bundles {
		bundleList = append(bundleList, b)
	}
	return r.registerBundles(bundleList)
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

	r.mu.Lock()
	for _, name := range names {
		r.bundles[name] = content
	}
	r.mu.Unlock()

	return nil
}

// DiscoverTagName loads a bundle and returns the custom element tag name it
// registers. It tries a fast regex pre-pass first; if that misses, it falls
// back to executing the bundle in a temporary QJS context.
func DiscoverTagName(bundle string) (string, error) {
	return discoverTagName(bundle)
}

// discoverTagNameFast extracts the tag name from a customElements.define()
// call using a regex, avoiding QJS entirely. Returns ("", false) when the
// regex cannot find a match.
func discoverTagNameFast(bundle string) (string, bool) {
	matches := defineRe.FindAllStringSubmatch(bundle, -1)
	if len(matches) == 0 {
		return "", false
	}
	return matches[len(matches)-1][1], true
}

func discoverTagName(bundle string) (string, error) {
	if tagName, ok := discoverTagNameFast(bundle); ok {
		return tagName, nil
	}

	engine, err := NewEngine()
	if err != nil {
		return "", err
	}
	defer engine.Close()

	return discoverTagNameEngine(engine, bundle)
}

// discoverTagNameEngine runs a single bundle through an already-initialized
// QJS engine and returns the last registered tag name.
func discoverTagNameEngine(engine *Engine, bundle string) (string, error) {
	if err := engine.LoadBundle(bundle); err != nil {
		return "", err
	}

	result, err := engine.ctx.Eval("discover.js", qjs.Code(discoveryJS))
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

	return names[len(names)-1], nil
}

// discoverTagNames runs a batch of bundles through a single reusable QJS
// engine, resetting between each. Returns a map of input index to tag name
// for bundles that successfully registered a custom element.
func discoverTagNames(bundles []string) (map[int]string, error) {
	if len(bundles) == 0 {
		return nil, nil
	}

	engine, err := NewEngine()
	if err != nil {
		return nil, err
	}
	defer engine.Close()

	results := make(map[int]string, len(bundles))

	for i, bundle := range bundles {
		name, err := discoverTagNameEngine(engine, bundle)
		if err == nil {
			results[i] = name
		}

		if i < len(bundles)-1 {
			if err := engine.Reset(); err != nil {
				return results, fmt.Errorf("resetting discovery engine: %w", err)
			}
		}
	}

	return results, nil
}
