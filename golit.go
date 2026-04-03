// Package golit provides a high-level API for server-side rendering Lit
// web components into Declarative Shadow DOM HTML.
//
// The Renderer type wraps the lower-level jsengine and transformer packages
// into a single entry point suitable for library consumers.
//
//	renderer, err := golit.NewRenderer(golit.RendererOptions{
//	    DefsDir: "bundles/",
//	})
//	if err != nil { log.Fatal(err) }
//	defer renderer.Close()
//
//	output, err := renderer.RenderFragment(`<my-el name="World"></my-el>`)
package golit

import (
	"github.com/sspriggs/golit/pkg/jsengine"
	"github.com/sspriggs/golit/pkg/transformer"
)

// RendererOptions configures a Renderer instance.
type RendererOptions struct {
	// DefsDir is a directory of pre-bundled .golit.bundle.js files.
	DefsDir string

	// SourcesDir is a directory of component .js/.ts source files to auto-bundle.
	SourcesDir string

	// ImportMap is a path to an import map JSON file.
	ImportMap string

	// Ignored is a list of custom element tag names to skip during SSR.
	Ignored []string

	// Preload is a list of module names to preload for dynamic imports.
	Preload []string
}

// Renderer holds a JS engine and bundle registry, providing methods
// to render Lit components into Declarative Shadow DOM.
type Renderer struct {
	engine   *jsengine.Engine
	registry *jsengine.Registry
	ignored  map[string]bool
}

// NewRenderer creates a Renderer, loading bundles and sources as specified.
func NewRenderer(opts RendererOptions) (*Renderer, error) {
	registry := jsengine.NewRegistry()

	if opts.DefsDir != "" {
		if err := registry.LoadDir(opts.DefsDir); err != nil {
			return nil, err
		}
	}

	if opts.SourcesDir != "" {
		if err := registry.LoadSourceDir(opts.SourcesDir); err != nil {
			return nil, err
		}
	}

	engine, err := jsengine.NewEngine()
	if err != nil {
		return nil, err
	}

	engine.SetPreloadModules(opts.Preload)

	ignored := make(map[string]bool, len(opts.Ignored))
	for _, tag := range opts.Ignored {
		ignored[tag] = true
	}

	return &Renderer{
		engine:   engine,
		registry: registry,
		ignored:  ignored,
	}, nil
}

// RenderHTML transforms a full HTML document, expanding custom elements
// into Declarative Shadow DOM. Reuses the Renderer's engine to avoid
// per-call engine creation overhead.
func (r *Renderer) RenderHTML(input string) (string, error) {
	return transformer.RenderHTMLWithEngine(input, r.engine, r.registry, r.ignored)
}

// RenderFragment transforms an HTML fragment, expanding custom elements
// into Declarative Shadow DOM. Reuses the Renderer's engine to avoid
// per-call engine creation overhead.
func (r *Renderer) RenderFragment(input string) (string, error) {
	return transformer.RenderFragmentWithEngine(input, r.engine, r.registry, r.ignored)
}

// TransformDir processes all HTML files in a directory tree using the
// Renderer's pre-loaded registry and ignored list.
func (r *Renderer) TransformDir(dir string, opts ...transformer.Options) (*transformer.Result, error) {
	o := transformer.Options{}
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Registry == nil {
		o.Registry = r.registry
	}
	if o.Ignored == nil {
		o.Ignored = r.ignored
	}
	return transformer.TransformDir(dir, o)
}

// RegisterComponent bundles inline JS/TS source and registers the
// resulting component for rendering. The source must call
// customElements.define() to register a tag name.
func (r *Renderer) RegisterComponent(source string) error {
	bundle, err := jsengine.BundleSource(source)
	if err != nil {
		return err
	}
	tagName, err := jsengine.DiscoverTagName(bundle)
	if err != nil {
		return err
	}
	r.registry.Register(tagName, bundle)
	return nil
}

// Registry returns the underlying bundle registry for direct manipulation.
func (r *Renderer) Registry() *jsengine.Registry {
	return r.registry
}

// Close releases the JS engine resources.
func (r *Renderer) Close() {
	r.engine.Close()
}
