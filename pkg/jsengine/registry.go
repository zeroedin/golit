package jsengine

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fastschema/qjs"
)

// defineRe extracts custom element tag names from customElements.define() calls
// without executing the module. Valid custom element names must contain a hyphen.
var defineRe = regexp.MustCompile(`customElements\s*\.\s*define\s*\(\s*['"]([a-z][a-z0-9]*(?:-[a-z0-9]+)+)['"]`)

// decoratorDefineRe matches the Lit @customElement decorator pattern used in
// thin ESM modules: customElement("tag-name") or customElement3("tag-name").
var decoratorDefineRe = regexp.MustCompile(`customElement\d*\s*\(\s*['"]([a-z][a-z0-9]*(?:-[a-z0-9]+)+)['"]`)


// Registry manages loaded component modules and tracks which tag names
// are available for rendering. All methods are safe for concurrent use.
type Registry struct {
	mu sync.RWMutex

	// modules maps tag names to their thin ES module JS content
	modules map[string]string

	// sharedRuntime is the shared runtime module source (loaded once per engine)
	sharedRuntime string

	// runtimeExternals lists the package prefixes bundled into the shared
	// runtime (e.g. "lit", "lit/*", "@rhds/tokens", "@rhds/tokens/*").
	// Used by Engine.shimDynamicImports to rewrite dynamic import() calls
	// for these packages to import("@golit/runtime").
	runtimeExternals []string

	// dynamicImportTargets lists specific module specifiers that appear in
	// dynamic import() calls in thin modules (e.g. "@rhds/tokens/css/default-theme.css.js").
	// These need to be bundled as standalone modules so the engine
	// can resolve them at runtime.
	dynamicImportTargets []string

	// dynamicModules maps specifier -> raw ESM source for modules that
	// need to be registered in QJS under their original specifier name
	// so dynamic import() calls resolve natively.
	dynamicModules map[string]string

	// baseDir is the directory from which modules were loaded, used to
	// resolve dynamic import targets relative to the correct node_modules.
	baseDir string

	// unregistered tracks custom element tags found but not in the registry
	unregistered map[string]bool

	// processedPaths tracks source file paths that have already been processed,
	// so discoverFromHTML can skip re-processing across multiple HTML files.
	processedPaths map[string]bool

	// bytecodeCache stores pre-compiled QJS bytecode keyed by module name.
	// Populated by the first engine during PreloadAll; subsequent engines
	// load from bytecode instead of re-parsing source.
	bytecodeCache map[string][]byte
}

// NewRegistry creates an empty registry for component modules and optional
// shared runtime state.
func NewRegistry() *Registry {
	return &Registry{
		modules:        make(map[string]string),
		unregistered:   make(map[string]bool),
		processedPaths: make(map[string]bool),
	}
}

// LoadDir loads all .golit.module.js files from a directory.
// If a _runtime.golit.module.js file is present, it is loaded as the shared runtime.
// Tag names are discovered via a regex pre-pass. Dynamic import() targets in
// the thin modules are collected and stored as dynamic import targets.
func (r *Registry) LoadDir(dir string) error {
	absDir, _ := filepath.Abs(dir)
	r.SetBaseDir(absDir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading modules directory %s: %w", dir, err)
	}

	moduleSources := make(map[string]string)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		fullPath := filepath.Join(dir, name)

		if name == "_runtime.golit.module.js" {
			data, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("reading runtime module %s: %w", name, err)
			}
			r.SetSharedRuntime(string(data))
			continue
		}

		if strings.HasSuffix(name, ".golit.module.js") {
			data, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("reading module %s: %w", name, err)
			}
			source := string(data)
			moduleSources[name] = source
			if tagName, ok := discoverTagNameFast(source); ok {
				r.RegisterModule(tagName, source)
			}
		}
	}

	if len(moduleSources) > 0 {
		targets := extractDynamicImportTargets(moduleSources)
		if len(targets) > 0 {
			r.SetDynamicImportTargets(targets)
			for _, target := range targets {
				modPath, err := ResolveModulePath(target, absDir)
				if err != nil {
					fmt.Fprintf(os.Stderr, "golit: warning: could not resolve dynamic module %s from %s: %v\n", target, absDir, err)
					continue
				}
				esm, err := BundleStandaloneModule(modPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "golit: warning: could not bundle dynamic module %s: %v\n", target, err)
					continue
				}
				r.SetDynamicModule(target, esm)
			}
		}

		unrewritten := ExtractUnrewrittenImports(moduleSources)
		for _, spec := range unrewritten {
			if r.dynamicModules != nil && r.dynamicModules[spec] != "" {
				continue
			}
			modPath, err := ResolveModulePath(spec, absDir)
			if err != nil {
				continue
			}
			esm, err := BundleStandaloneModule(modPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "golit: warning: could not bundle standalone module %s: %v\n", spec, err)
				continue
			}
			r.SetDynamicModule(spec, esm)
		}
	}

	return nil
}

// LoadFile loads a single .golit.module.js file and discovers its tag name.
// Modules that don't register any custom elements are silently skipped.
func (r *Registry) LoadFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading module %s: %w", path, err)
	}

	source := string(data)
	tagName, ok := discoverTagNameFast(source)
	if !ok {
		return nil
	}

	r.Register(tagName, source)
	return nil
}

// Register adds a module by tag name directly (for programmatic use).
func (r *Registry) Register(tagName string, source string) {
	r.mu.Lock()
	r.modules[tagName] = source
	r.mu.Unlock()
}

// RegisterModule adds a thin ES module by tag name (alias for Register).
func (r *Registry) RegisterModule(tagName string, source string) {
	r.Register(tagName, source)
}

// Lookup returns the ES module JS for a given tag name, or "" if not found.
func (r *Registry) Lookup(tagName string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.modules[tagName]
}

// LookupModule returns the ES module JS for a given tag name, or "" if not found.
// Alias for Lookup.
func (r *Registry) LookupModule(tagName string) string {
	return r.Lookup(tagName)
}

// SharedRuntime returns the shared runtime module source, or "" if not set.
func (r *Registry) SharedRuntime() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.sharedRuntime
}

// SetSharedRuntime sets the shared runtime module source.
func (r *Registry) SetSharedRuntime(source string) {
	r.mu.Lock()
	r.sharedRuntime = source
	r.mu.Unlock()
}

// RuntimeExternals returns a copy of the package prefixes bundled into the shared runtime.
func (r *Registry) RuntimeExternals() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]string(nil), r.runtimeExternals...)
}

// SetRuntimeExternals stores the package prefixes bundled into the shared runtime.
func (r *Registry) SetRuntimeExternals(externals []string) {
	r.mu.Lock()
	r.runtimeExternals = externals
	r.mu.Unlock()
}

// DynamicImportTargets returns a copy of the module specifiers found in
// dynamic import() calls within thin modules.
func (r *Registry) DynamicImportTargets() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]string(nil), r.dynamicImportTargets...)
}

// SetDynamicImportTargets stores the dynamic import target specifiers.
func (r *Registry) SetDynamicImportTargets(targets []string) {
	r.mu.Lock()
	r.dynamicImportTargets = targets
	r.mu.Unlock()
}

// DynamicModules returns a copy of the specifier -> ESM source map for
// modules that should be registered as named QJS modules.
func (r *Registry) DynamicModules() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.dynamicModules == nil {
		return nil
	}
	cp := make(map[string]string, len(r.dynamicModules))
	for k, v := range r.dynamicModules {
		cp[k] = v
	}
	return cp
}

// SetDynamicModule registers a standalone ES module under its specifier name.
func (r *Registry) SetDynamicModule(specifier, source string) {
	r.mu.Lock()
	if r.dynamicModules == nil {
		r.dynamicModules = make(map[string]string)
	}
	r.dynamicModules[specifier] = source
	r.mu.Unlock()
}

// BaseDir returns the directory from which modules were loaded.
func (r *Registry) BaseDir() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.baseDir
}

// SetBaseDir stores the base directory for module resolution.
func (r *Registry) SetBaseDir(dir string) {
	r.mu.Lock()
	r.baseDir = dir
	r.mu.Unlock()
}

// Has returns true if a module is registered for the given tag name.
func (r *Registry) Has(tagName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.modules[tagName]
	return ok
}

// TagNames returns all registered tag names.
func (r *Registry) TagNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}

// SetBytecode stores pre-compiled bytecode for a module name.
func (r *Registry) SetBytecode(name string, bc []byte) {
	r.mu.Lock()
	if r.bytecodeCache == nil {
		r.bytecodeCache = make(map[string][]byte)
	}
	r.bytecodeCache[name] = bc
	r.mu.Unlock()
}

// LookupBytecode returns pre-compiled bytecode for a module name, or nil.
func (r *Registry) LookupBytecode(name string) []byte {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.bytecodeCache == nil {
		return nil
	}
	return r.bytecodeCache[name]
}

// HasBytecode returns true if any bytecode has been cached.
func (r *Registry) HasBytecode() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.bytecodeCache) > 0
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

// LoadSourceDir bundles all .js/.ts files in a directory tree as thin ES modules
// and registers them. Also builds and registers the shared runtime.
func (r *Registry) LoadSourceDir(dir string) error {
	absDir, _ := filepath.Abs(dir)
	r.SetBaseDir(absDir)

	var paths []string
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: skipping %s: %v\n", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasSuffix(name, ".d.ts") || strings.HasSuffix(name, ".golit.module.js") {
			return nil
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

	nodeModulesDir := FindNodeModules(paths[0])

	externals, err := DiscoverExternalPackages(paths, nodeModulesDir)
	if err != nil {
		return fmt.Errorf("discovering external packages: %w", err)
	}

	modules, err := BundleComponentModules(paths, BundleOptions{
		ExternalPackages: externals,
	})
	if err != nil {
		return fmt.Errorf("batch bundling sources: %w", err)
	}

	if r.SharedRuntime() == "" && nodeModulesDir != "" {
		rt, err := BundleSharedRuntime(nodeModulesDir, modules)
		if err != nil {
			return fmt.Errorf("building shared runtime: %w", err)
		}
		r.SetSharedRuntime(rt)
	}

	r.SetRuntimeExternals(externals)

	targets := extractDynamicImportTargets(modules)
	if len(targets) > 0 {
		r.SetDynamicImportTargets(targets)
	}

	modules = RewriteModuleImports(modules, externals)

	// After rewriting, find default imports that kept their original specifier
	// (not rewritten to @golit/runtime) and bundle them as standalone modules.
	unrewritten := ExtractUnrewrittenImports(modules)
	for _, spec := range unrewritten {
		modPath, err := ResolveModulePath(spec, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: could not resolve unrewritten import %s: %v\n", spec, err)
			continue
		}
		esm, err := BundleStandaloneModule(modPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: could not bundle standalone module %s: %v\n", spec, err)
			continue
		}
		r.SetDynamicModule(spec, esm)
	}

	for _, source := range modules {
		if tagName, ok := discoverTagNameFast(source); ok {
			r.RegisterModule(tagName, source)
		}
	}

	return nil
}

// LoadCompiled loads a single pre-compiled .golit.compiled.js artifact
// that contains all components and a __golitRegistry manifest.
// The compiled artifact is treated as a single module that registers
// all components when evaluated.
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
		r.modules[name] = content
	}
	r.mu.Unlock()

	return nil
}

// DiscoverTagName extracts the custom element tag name from a module source
// using a regex. Returns the tag name or an error if none found.
func DiscoverTagName(source string) (string, error) {
	if tagName, ok := discoverTagNameFast(source); ok {
		return tagName, nil
	}
	return "", fmt.Errorf("no custom element registration found in source")
}

// discoverTagNameFast extracts the tag name from a customElements.define()
// or @customElement() decorator call using a regex, avoiding QJS entirely.
// Returns ("", false) when the regex cannot find a match.
func discoverTagNameFast(source string) (string, bool) {
	idx := strings.LastIndex(source, "customElements")
	if idx >= 0 {
		if match := defineRe.FindStringSubmatch(source[idx:]); match != nil {
			return match[1], true
		}
		matches := defineRe.FindAllStringSubmatch(source, -1)
		if len(matches) > 0 {
			return matches[len(matches)-1][1], true
		}
	}

	matches := decoratorDefineRe.FindAllStringSubmatch(source, -1)
	if len(matches) > 0 {
		return matches[len(matches)-1][1], true
	}

	return "", false
}
