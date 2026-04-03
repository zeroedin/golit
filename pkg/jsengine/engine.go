package jsengine

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
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

// Engine executes Lit components in QJS for SSR rendering.
type Engine struct {
	runtime        *qjs.Runtime
	ctx            *qjs.Context
	loaded         map[string]bool // track which bundles have been loaded
	preloadModules []string        // module names available via __preloadedModules
}

// NewEngine creates a new QJS engine instance.
func NewEngine() (*Engine, error) {
	rt, err := qjs.New()
	if err != nil {
		return nil, fmt.Errorf("creating QJS runtime: %w", err)
	}
	e := &Engine{
		runtime: rt,
		ctx:     rt.Context(),
		loaded:  make(map[string]bool),
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
	_, err := e.ctx.Eval("helpers.js", qjs.Code(helpersJS))
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

// LoadBundle loads a pre-bundled component JS file into the engine.
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

// shimDynamicImports replaces dynamic import("module") expressions with
// Promise.resolve(globalThis.__preloadedModules["module"]) for preloaded modules.
func (e *Engine) shimDynamicImports(code string) string {
	for _, mod := range e.preloadModules {
		// Replace: import("prism-esm") -> Promise.resolve(globalThis.__preloadedModules["prism-esm"])
		// Also handle: import('prism-esm') with single quotes
		// And subpath imports: import("prism-esm/components/prism-css.js")
		dq := `import("` + mod
		sq := `import('` + mod
		replacement := `Promise.resolve(globalThis.__preloadedModules["` + mod + `"] || {})/*golit-shimmed:`
		// Replace double-quoted imports
		for strings.Contains(code, dq) {
			// Find the full import expression: import("module...") 
			idx := strings.Index(code, dq)
			// Find closing quote + paren
			end := strings.Index(code[idx+len(dq):], `")`)
			if end < 0 {
				break
			}
			full := code[idx : idx+len(dq)+end+2]
			code = strings.Replace(code, full, replacement+full+`*/`, 1)
		}
		// Replace single-quoted imports
		for strings.Contains(code, sq) {
			idx := strings.Index(code, sq)
			end := strings.Index(code[idx+len(sq):], `')`)
			if end < 0 {
				break
			}
			full := code[idx : idx+len(sq)+end+2]
			code = strings.Replace(code, full, replacement+full+`*/`, 1)
		}
	}
	return code
}

// LoadBundleForTag loads a bundle from the registry for a specific tag name.
// Returns false if the tag is not in the registry.
func (e *Engine) LoadBundleForTag(tagName string, registry *Registry) bool {
	if e.loaded[tagName] {
		return true
	}

	bundle := registry.Lookup(tagName)
	if bundle == "" {
		return false
	}

	if err := e.LoadBundle(bundle); err != nil {
		return false
	}

	e.loaded[tagName] = true
	return true
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
					const result = el.render();
					html = __collectTemplateResult(result);
				}

				let css = '';
				if (Ctor.styles) {
					css = extractStyles(Ctor.styles);
				} else if (Ctor.elementStyles) {
					css = extractStyles(Ctor.elementStyles);
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

// RenderBatch renders multiple custom elements in a single QJS Eval call,
// reducing Go-to-JS boundary crossings. Each element is still rendered
// individually within QJS; the batching is at the transport layer only.
func (e *Engine) RenderBatch(requests []BatchRequest) ([]BatchResult, error) {
	if len(requests) == 0 {
		return nil, nil
	}

	reqJSON, err := json.Marshal(requests)
	if err != nil {
		return nil, fmt.Errorf("marshaling batch requests: %w", err)
	}

	script := fmt.Sprintf(`
		(function() {
			const requests = %s;
			const results = [];
			for (const req of requests) {
				const Ctor = customElements.get(req.tagName);
				if (!Ctor) {
					results.push({ id: req.id, error: 'Element <' + req.tagName + '> not registered' });
					continue;
				}
				try {
					const el = new Ctor();
					for (const [key, value] of Object.entries(req.attrs || {})) {
						el.setAttribute(key, value);
						const propName = attributeToProperty(Ctor, key);
						if (propName) {
							const propConfig = getPropertyConfig(Ctor, propName);
							el[propName] = coerceValue(value, propConfig);
						}
					}
					let html = '';
					if (typeof el.render === 'function') {
						html = __collectTemplateResult(el.render());
					}
					let css = '';
					if (Ctor.styles) { css = extractStyles(Ctor.styles); }
					else if (Ctor.elementStyles) { css = extractStyles(Ctor.elementStyles); }
					results.push({ id: req.id, html: html, css: css, tagName: req.tagName });
				} catch(e) {
					results.push({ id: req.id, error: e.message, tagName: req.tagName });
				}
			}
			return JSON.stringify(results);
		})();
	`, string(reqJSON))

	result, err := e.ctx.Eval("render-batch.js", qjs.Code(script))
	if err != nil {
		return nil, fmt.Errorf("batch render eval: %w", err)
	}

	var results []BatchResult
	if err := json.Unmarshal([]byte(result.String()), &results); err != nil {
		return nil, fmt.Errorf("parsing batch results: %w (raw: %s)", err, result.String())
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
