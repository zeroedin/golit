package jsengine

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/fastschema/qjs"
)

// validTagName matches the custom element name spec: starts with [a-z],
// contains at least one hyphen, and only [a-z0-9-._] characters.
var validTagName = regexp.MustCompile(`^[a-z][a-z0-9]*(?:[-._][a-z0-9]+)+$`)

//go:embed helpers.js
var helpersJS string

// RenderResult contains the output of rendering a custom element.
type RenderResult struct {
	// HTML is the rendered shadow DOM content (the template output).
	HTML string

	// CSS is the component's scoped styles (from static styles).
	CSS string

	// TagName is the confirmed custom element tag name.
	TagName string
}

// renderCacheEntry stores the cached result of rendering a component
// with a specific set of attributes.
type renderCacheEntry struct {
	html string
	css  string
}

// Engine executes Lit components in QJS for SSR rendering.
type Engine struct {
	runtime          *qjs.Runtime
	ctx              *qjs.Context
	loaded           map[string]bool             // track which bundles have been loaded
	preloadModules   []string                    // module names available via __preloadedModules
	runtimeExternals []string                    // package prefixes bundled into @golit/runtime
	renderCache      map[string]renderCacheEntry // L1: per-engine, lock-free
	sharedCache      *sharedRenderCache          // L2: shared across pool engines
	renderFnReady    bool                        // whether __golitRenderBatch is registered
}

// NewEngine creates a new QJS engine instance.
func NewEngine() (*Engine, error) {
	rt, err := qjs.New()
	if err != nil {
		return nil, fmt.Errorf("creating QJS runtime: %w", err)
	}
	e := &Engine{
		runtime:     rt,
		ctx:         rt.Context(),
		loaded:      make(map[string]bool),
		renderCache: make(map[string]renderCacheEntry),
	}
	if err := e.initHelpers(); err != nil {
		rt.Close()
		return nil, fmt.Errorf("initializing helpers: %w", err)
	}
	return e, nil
}

// initHelpers loads the SSR rendering helper functions into the QJS
// global scope once, so they don't need to be re-sent with every
// RenderElement call.
func (e *Engine) initHelpers() error {
	if err := e.injectSSRStringGlobals(); err != nil {
		return fmt.Errorf("inject SSR globals: %w", err)
	}
	e.installFetchBridge()
	if _, err := e.ctx.Eval("helpers.js", qjs.Code(helpersJS)); err != nil {
		return err
	}
	_, err := e.ctx.Eval("css-cache-init.js", qjs.Code(
		`if(!globalThis.__cssCache)globalThis.__cssCache=new Map();`))
	return err
}

// Close releases the QJS runtime resources.
func (e *Engine) Close() {
	e.runtime.Close()
}

// Reset tears down the current QJS runtime and creates a new one,
// clearing all global state. Previously loaded bundles must be
// re-loaded afterward.
func (e *Engine) Reset() error {
	e.runtime.Close()
	rt, err := qjs.New()
	if err != nil {
		return fmt.Errorf("resetting QJS runtime: %w", err)
	}
	e.runtime = rt
	e.ctx = rt.Context()
	e.loaded = make(map[string]bool)
	e.renderCache = make(map[string]renderCacheEntry)
	e.renderFnReady = false
	if err := e.initHelpers(); err != nil {
		return fmt.Errorf("re-initializing helpers: %w", err)
	}
	return nil
}

// SetPreloadModules tells the engine which module names have been preloaded
// into globalThis.__preloadedModules. Dynamic import() calls for these modules
// will be shimmed to resolve from the preloaded registry.
func (e *Engine) SetPreloadModules(modules []string) {
	e.preloadModules = modules
}

// SetRuntimeExternals tells the engine which package prefixes are bundled
// into the @golit/runtime module. Dynamic import() calls for these packages
// will be rewritten to import("@golit/runtime").
func (e *Engine) SetRuntimeExternals(externals []string) {
	e.runtimeExternals = externals
}

// LoadBundle loads a pre-bundled component JS file into the engine as a script.
// Can be called multiple times to load multiple components.
// If preload modules are set, dynamic import() calls for those modules
// are replaced with synchronous lookups into __preloadedModules.
func (e *Engine) LoadBundle(bundle string) error {
	code := e.shimDynamicImports(bundle)
	_, err := e.ctx.Eval("bundle.js", qjs.Code(code))
	if err != nil {
		return fmt.Errorf("loading bundle: %w", err)
	}
	return nil
}

// LoadModule pre-registers a named ES module in the QJS module cache
// without executing it. Other modules can then import from it by name.
func (e *Engine) LoadModule(name string, source string) error {
	_, err := e.runtime.Load(name, qjs.Code(source))
	if err != nil {
		return fmt.Errorf("loading module %s: %w", name, err)
	}
	return nil
}

// EvalModule evaluates an ES module, executing its body and registering
// any custom elements it defines. Use for component modules that need
// their side-effects (like customElements.define) to run.
func (e *Engine) EvalModule(name string, source string) error {
	code := e.shimDynamicImports(source)
	_, err := e.ctx.Eval(name, qjs.Code(code), qjs.TypeModule())
	if err != nil {
		return fmt.Errorf("evaluating module %s: %w", name, err)
	}
	return nil
}

// CompileModule compiles an ES module to bytecode without executing it.
func (e *Engine) CompileModule(name string, source string) ([]byte, error) {
	code := e.shimDynamicImports(source)
	bc, err := e.ctx.Compile(name, qjs.Code(code), qjs.TypeModule(), qjs.FlagCompileOnly())
	if err != nil {
		return nil, fmt.Errorf("compiling module %s: %w", name, err)
	}
	return bc, nil
}

// EvalModuleBytecode loads and executes a pre-compiled ES module from bytecode.
func (e *Engine) EvalModuleBytecode(name string, bc []byte) error {
	_, err := e.ctx.Eval(name, qjs.Bytecode(bc))
	if err != nil {
		return fmt.Errorf("evaluating bytecode module %s: %w", name, err)
	}
	return nil
}

// LoadModuleBytecode loads a pre-compiled module into the module cache
// from bytecode without executing it.
func (e *Engine) LoadModuleBytecode(name string, bc []byte) error {
	_, err := e.runtime.Load(name, qjs.Bytecode(bc))
	if err != nil {
		return fmt.Errorf("loading bytecode module %s: %w", name, err)
	}
	return nil
}

// shimPattern holds a pre-computed (module, quote-style) pair so the
// scan loop in shimDynamicImports avoids per-iteration allocations.
type shimPattern struct {
	prefix string // e.g. `import("prism-esm`
	close  string // e.g. `")`
	mod    string
}

// shimDynamicImports replaces dynamic import("module") expressions with
// Promise.resolve(globalThis.__preloadedModules["module"]) for preloaded modules,
// and rewrites import("pkg/...") to import("@golit/runtime") for packages
// bundled into the shared runtime.
// Handles both quote styles and subpath imports (e.g. import("mod/sub.js")).
func (e *Engine) shimDynamicImports(code string) string {
	if (len(e.preloadModules) == 0 && len(e.runtimeExternals) == 0) || !strings.Contains(code, "import(") {
		return code
	}

	patterns := make([]shimPattern, 0, len(e.preloadModules)*2)
	for _, mod := range e.preloadModules {
		patterns = append(patterns,
			shimPattern{prefix: `import("` + mod, close: `")`, mod: mod},
			shimPattern{prefix: `import('` + mod, close: `')`, mod: mod},
		)
	}

	const importOpen = "import("
	var b strings.Builder
	b.Grow(len(code) + len(code)/10)

	pos := 0
	for {
		idx := strings.Index(code[pos:], importOpen)
		if idx < 0 {
			b.WriteString(code[pos:])
			break
		}

		b.WriteString(code[pos : pos+idx])
		matchStart := pos + idx
		matched := false

		for _, p := range patterns {
			if matchStart+len(p.prefix) > len(code) ||
				code[matchStart:matchStart+len(p.prefix)] != p.prefix {
				continue
			}
			nextChar := code[matchStart+len(p.prefix)]
			if nextChar != p.close[0] && nextChar != '/' {
				continue
			}
			end := strings.Index(code[matchStart+len(p.prefix):], p.close)
			if end < 0 {
				continue
			}
			full := code[matchStart : matchStart+len(p.prefix)+end+len(p.close)]
			b.WriteString(`Promise.resolve(globalThis.__preloadedModules["`)
			b.WriteString(p.mod)
			b.WriteString(`"] || {})/*golit-shimmed:`)
			b.WriteString(full)
			b.WriteString(`*/`)
			pos = matchStart + len(full)
			matched = true
			break
		}

		if !matched && len(e.runtimeExternals) > 0 {
			if rewritten, newPos := e.shimRuntimeImport(code, matchStart); rewritten != "" {
				b.WriteString(rewritten)
				pos = newPos
				matched = true
			}
		}

		if !matched {
			b.WriteString(importOpen)
			pos = matchStart + len(importOpen)
		}
	}

	return b.String()
}

// shimRuntimeImport checks if a dynamic import() at the given position
// matches a runtime external package and rewrites it to import("@golit/runtime").
// Skips specifiers that have a standalone module registered (e.loaded),
// allowing QJS native module resolution to handle them instead.
func (e *Engine) shimRuntimeImport(code string, matchStart int) (string, int) {
	for _, quote := range []byte{'"', '\''} {
		prefix := `import(` + string(quote)
		if matchStart+len(prefix) > len(code) ||
			code[matchStart:matchStart+len(prefix)] != prefix {
			continue
		}
		closeStr := string(quote) + ")"
		end := strings.Index(code[matchStart+len(prefix):], closeStr)
		if end < 0 {
			continue
		}
		specifier := code[matchStart+len(prefix) : matchStart+len(prefix)+end]
		if e.loaded[specifier] {
			return "", 0
		}
		if matchesExternals(specifier, e.runtimeExternals) {
			full := code[matchStart : matchStart+len(prefix)+end+len(closeStr)]
			rewritten := `import(` + string(quote) + `@golit/runtime` + string(quote) + `)/*golit-runtime:` + full + `*/`
			return rewritten, matchStart + len(full)
		}
	}
	return "", 0
}

// LoadBundleForTag loads a component from the registry for a specific tag name.
// Returns (true, nil) if loaded or already loaded, (false, nil) if the tag
// is not in the registry, or (false, err) if loading failed.
// If a shared runtime is present, it is loaded first. All registry content
// is evaluated as ES modules via EvalModule (registry stores .golit.module.js).
func (e *Engine) LoadBundleForTag(tagName string, registry *Registry) (bool, error) {
	if e.loaded[tagName] {
		return true, nil
	}

	source := registry.Lookup(tagName)
	if source == "" {
		return false, nil
	}

	if rt := registry.SharedRuntime(); rt != "" && !e.loaded["@golit/runtime"] {
		if err := e.LoadModule("@golit/runtime", rt); err != nil {
			return false, fmt.Errorf("loading shared runtime: %w", err)
		}
		e.loaded["@golit/runtime"] = true
	}

	if ext := registry.RuntimeExternals(); len(ext) > 0 && len(e.runtimeExternals) == 0 {
		e.runtimeExternals = ext
	}

	if !e.loaded["@golit/dynamic-modules"] {
		for specifier, modSource := range registry.DynamicModules() {
			if !e.loaded[specifier] {
				if err := e.LoadModule(specifier, modSource); err != nil {
					return false, fmt.Errorf("loading dynamic module %s for <%s>: %w", specifier, tagName, err)
				}
				e.loaded[specifier] = true
			}
		}
		e.loaded["@golit/dynamic-modules"] = true
	}

	if err := e.EvalModule(tagName+".js", source); err != nil {
		return false, fmt.Errorf("loading module for <%s>: %w", tagName, err)
	}

	e.loaded[tagName] = true
	return true, nil
}

// RenderElement creates an instance of a custom element, sets attributes
// on it, calls render(), and returns the rendered HTML and CSS.
func (e *Engine) RenderElement(tagName string, attrs map[string]string) (*RenderResult, error) {
	if !validTagName.MatchString(tagName) {
		return nil, fmt.Errorf("invalid custom element name: %q", tagName)
	}

	attrsJSON, err := json.Marshal(attrs)
	if err != nil {
		return nil, fmt.Errorf("marshaling attrs: %w", err)
	}

	script := fmt.Sprintf(`
		(function() {
			const Ctor = customElements.get('%s');
			if (!Ctor) {
				return JSON.stringify({ error: 'Element <%s> not registered' });
			}

			try {
				const el = new Ctor();
				const attrs = %s;

				for (const [key, value] of Object.entries(attrs)) {
					el.setAttribute(key, value);
					const propName = attributeToProperty(Ctor, key);
					if (propName) {
						const propConfig = getPropertyConfig(Ctor, propName);
						el[propName] = coerceValue(value, propConfig);
					}
				}

				let html = '';
				if (typeof el.render === 'function') {
					html = __collectTemplateResult(el.render(), true);
				} else {
					html = __collectTemplateResult(null, true);
				}

				let css = '';
				if (globalThis.__cssCache.has(Ctor)) {
					css = globalThis.__cssCache.get(Ctor);
				} else {
					if (Ctor.styles) { css = extractStyles(Ctor.styles); }
					else if (Ctor.elementStyles) { css = extractStyles(Ctor.elementStyles); }
					globalThis.__cssCache.set(Ctor, css);
				}

				return JSON.stringify({ html, css, tagName: '%s' });
			} catch(e) {
				return JSON.stringify({ error: e.message, stack: e.stack || '' });
			}
		})();
	`, tagName, tagName, string(attrsJSON), tagName)

	result, err := e.ctx.Eval("render.js", qjs.Code(script))
	if err != nil {
		return nil, fmt.Errorf("rendering <%s>: %w", tagName, err)
	}

	var output struct {
		HTML    string `json:"html"`
		CSS     string `json:"css"`
		TagName string `json:"tagName"`
		Error   string `json:"error"`
		Stack   string `json:"stack"`
	}

	if err := json.Unmarshal([]byte(result.String()), &output); err != nil {
		return nil, fmt.Errorf("parsing render result: %w (raw: %s)", err, result.String())
	}

	if output.Error != "" {
		return nil, fmt.Errorf("JS error rendering <%s>: %s\n%s", tagName, output.Error, output.Stack)
	}

	return &RenderResult{
		HTML:    output.HTML,
		CSS:     strings.TrimSpace(output.CSS),
		TagName: output.TagName,
	}, nil
}

// BatchRequest describes a single element to render in a batch call.
type BatchRequest struct {
	ID      int               `json:"id"`
	TagName string            `json:"tagName"`
	Attrs   map[string]string `json:"attrs"`
}

// BatchResult contains the output of rendering a single element in a batch.
type BatchResult struct {
	ID      int    `json:"id"`
	HTML    string `json:"html"`
	CSS     string `json:"css"`
	TagName string `json:"tagName"`
	Error   string `json:"error,omitempty"`
}

// renderCacheKey builds a deterministic cache key from tag name and attributes.
func renderCacheKey(tagName string, attrs map[string]string) string {
	if len(attrs) == 0 {
		return tagName
	}
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(tagName)
	for _, k := range keys {
		b.WriteByte(0)
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(attrs[k])
	}
	return b.String()
}

// registerRenderFn registers the batch render function in JS once,
// so subsequent calls avoid re-parsing the render loop code.
func (e *Engine) registerRenderFn() error {
	if e.renderFnReady {
		return nil
	}
	const renderFnJS = `globalThis.__golitRenderBatch=function(requestsJSON){` +
		`const requests=JSON.parse(requestsJSON);const results=[];` +
		`for(const req of requests){` +
		`const Ctor=customElements.get(req.tagName);` +
		`if(!Ctor){results.push({id:req.id,error:'Element <'+req.tagName+'> not registered',tagName:req.tagName});continue;}` +
		`try{` +
		`const el=new Ctor();` +
		`for(const [key,value] of Object.entries(req.attrs||{})){` +
		`el.setAttribute(key,value);` +
		`const propName=attributeToProperty(Ctor,key);` +
		`if(propName){const propConfig=getPropertyConfig(Ctor,propName);el[propName]=coerceValue(value,propConfig);}}` +
		`let html='';if(typeof el.render==='function'){html=__collectTemplateResult(el.render(),true);}else{html=__collectTemplateResult(null,true);}` +
		`let css='';if(globalThis.__cssCache.has(Ctor)){css=globalThis.__cssCache.get(Ctor);}` +
		`else{if(Ctor.styles){css=extractStyles(Ctor.styles);}` +
		`else if(Ctor.elementStyles){css=extractStyles(Ctor.elementStyles);}` +
		`globalThis.__cssCache.set(Ctor,css);}` +
		`results.push({id:req.id,html:html,css:css,tagName:req.tagName});` +
		`}catch(e){results.push({id:req.id,error:String(e&&e.message||e),tagName:req.tagName});}}` +
		`return JSON.stringify(results);};`
	_, err := e.ctx.Eval("golit-render-fn.js", qjs.Code(renderFnJS))
	if err != nil {
		return fmt.Errorf("registering render function: %w", err)
	}
	e.renderFnReady = true
	return nil
}

// RenderBatch renders multiple custom elements in a single QJS call.
// Uses Go-side render result caching to skip JS execution for
// previously-seen (tagName, attrs) combinations, and calls a
// pre-registered JS function to avoid re-parsing the render loop.
func (e *Engine) RenderBatch(requests []BatchRequest) ([]BatchResult, error) {
	if len(requests) == 0 {
		return nil, nil
	}

	results := make([]BatchResult, len(requests))
	var uncached []BatchRequest
	uncachedIdx := make([]int, 0, len(requests))
	var uncachedKeys []string

	for i, req := range requests {
		key := renderCacheKey(req.TagName, req.Attrs)
		if entry, ok := e.renderCache[key]; ok {
			results[i] = BatchResult{
				ID:      req.ID,
				HTML:    entry.html,
				CSS:     entry.css,
				TagName: req.TagName,
			}
		} else if e.sharedCache != nil {
			if entry, ok := e.sharedCache.get(key); ok {
				e.renderCache[key] = entry
				results[i] = BatchResult{
					ID:      req.ID,
					HTML:    entry.html,
					CSS:     entry.css,
					TagName: req.TagName,
				}
			} else {
				uncached = append(uncached, req)
				uncachedIdx = append(uncachedIdx, i)
				uncachedKeys = append(uncachedKeys, key)
			}
		} else {
			uncached = append(uncached, req)
			uncachedIdx = append(uncachedIdx, i)
			uncachedKeys = append(uncachedKeys, key)
		}
	}

	if len(uncached) == 0 {
		return results, nil
	}

	if err := e.registerRenderFn(); err != nil {
		return nil, err
	}

	reqJSON, err := json.Marshal(uncached)
	if err != nil {
		return nil, fmt.Errorf("marshaling batch requests: %w", err)
	}

	var script strings.Builder
	script.Grow(len(reqJSON) + 40)
	script.WriteString("__golitRenderBatch('")
	for _, b := range reqJSON {
		switch b {
		case '\'':
			script.WriteString("\\'")
		case '\\':
			script.WriteString("\\\\")
		case '\n':
			script.WriteString("\\n")
		case '\r':
			script.WriteString("\\r")
		default:
			script.WriteByte(b)
		}
	}
	script.WriteString("')")

	result, err := e.ctx.Eval("render-batch.js", qjs.Code(script.String()))
	if err != nil {
		return nil, fmt.Errorf("batch render eval: %w", err)
	}

	var jsResults []BatchResult
	if err := json.Unmarshal([]byte(result.String()), &jsResults); err != nil {
		return nil, fmt.Errorf("parsing batch results: %w (raw: %s)", err, result.String())
	}

	jsResultMap := make(map[int]BatchResult, len(jsResults))
	for _, r := range jsResults {
		jsResultMap[r.ID] = r
	}

	for j, origIdx := range uncachedIdx {
		req := uncached[j]
		r, ok := jsResultMap[req.ID]
		if !ok {
			results[origIdx] = BatchResult{
				ID:      req.ID,
				TagName: req.TagName,
				Error:   fmt.Sprintf("no result returned from JS for <%s>", req.TagName),
			}
			continue
		}
		results[origIdx] = r

		if r.Error == "" {
			entry := renderCacheEntry{html: r.HTML, css: r.CSS}
			e.renderCache[uncachedKeys[j]] = entry
			if e.sharedCache != nil {
				e.sharedCache.set(uncachedKeys[j], entry)
			}
		}
	}

	return results, nil
}

// IsRegistered checks if a tag name is registered in the QJS custom elements registry.
func (e *Engine) IsRegistered(tagName string) bool {
	if !validTagName.MatchString(tagName) {
		return false
	}
	result, err := e.ctx.Eval("check.js", qjs.Code(fmt.Sprintf(
		`customElements.get('%s') !== undefined ? 'true' : 'false'`, tagName)))
	if err != nil {
		return false
	}
	return result.String() == "true"
}
