package jsengine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fastschema/qjs"
)

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
	return &Engine{
		runtime: rt,
		ctx:     rt.Context(),
		loaded:  make(map[string]bool),
	}, nil
}

// Close releases the QJS runtime resources.
func (e *Engine) Close() {
	e.runtime.Close()
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

				// Set attributes AND properties on the element
				for (const [key, value] of Object.entries(attrs)) {
					el.setAttribute(key, value);
					const propName = attributeToProperty(Ctor, key);
					if (propName) {
						const propConfig = getPropertyConfig(Ctor, propName);
						el[propName] = coerceValue(value, propConfig);
					}
				}

				// Get rendered HTML
				let html = '';
				if (typeof el.render === 'function') {
					const result = el.render();
					html = __collectTemplateResult(result);
				}

				// Extract CSS from static styles
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

		function attributeToProperty(Ctor, attrName) {
			if (Ctor.__attributeToPropertyMap) {
				return Ctor.__attributeToPropertyMap.get(attrName);
			}
			if (Ctor.elementProperties) {
				for (const [propName, config] of Ctor.elementProperties) {
					const mappedAttr = (config && config.attribute !== undefined)
						? (config.attribute === false ? null : config.attribute)
						: propName.toLowerCase();
					if (mappedAttr === attrName) return propName;
				}
			}
			return attrName;
		}

		function getPropertyConfig(Ctor, propName) {
			if (Ctor.elementProperties) {
				return Ctor.elementProperties.get(propName) || {};
			}
			return {};
		}

		function coerceValue(value, config) {
			const type = config && config.type;
			if (type === Number) return Number(value);
			if (type === Boolean) return value !== 'false';
			return value;
		}

		function extractStyles(styles) {
			if (!styles) return '';
			if (typeof styles === 'string') return styles;
			if (Array.isArray(styles)) {
				return styles.map(s => extractStyles(s)).filter(Boolean).join('\n');
			}
			// CSSResult object from Lit's css tagged template
			if (styles.cssText !== undefined) return styles.cssText;
			// Adopted stylesheet
			if (styles._$cssResult$) return styles.cssText || '';
			return '';
		}
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

// IsRegistered checks if a tag name is registered in the QJS custom elements registry.
func (e *Engine) IsRegistered(tagName string) bool {
	result, err := e.ctx.Eval("check.js", qjs.Code(fmt.Sprintf(
		`customElements.get('%s') !== undefined ? 'true' : 'false'`, tagName)))
	if err != nil {
		return false
	}
	return result.String() == "true"
}
