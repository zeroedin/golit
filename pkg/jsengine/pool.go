package jsengine

import (
	"fmt"
	"os"
)

// EnginePool manages a fixed set of QJS Engine instances for concurrent use.
// Each engine is an isolated WASM module instance (via wazero) and can safely
// run in its own goroutine.
type EnginePool struct {
	engines chan *Engine
	size    int
}

// NewEnginePool creates a pool of size engines. Each engine is initialized
// with helpers but no bundles loaded yet — call PreloadAll before dispatch.
func NewEnginePool(size int) (*EnginePool, error) {
	if size < 1 {
		size = 1
	}
	p := &EnginePool{
		engines: make(chan *Engine, size),
		size:    size,
	}
	for i := 0; i < size; i++ {
		e, err := NewEngine()
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("creating engine %d/%d: %w", i+1, size, err)
		}
		p.engines <- e
	}
	return p, nil
}

// PreloadAll drains every engine from the pool, configures preload modules,
// loads any raw preload bundles, bundles dynamic import targets as preloads,
// loads the shared runtime (if present), then loads all registry component
// bundles/modules.
// After this call the registry must be treated as read-only.
func (p *EnginePool) PreloadAll(registry *Registry, preloadModules []string, preloadBundles ...string) error {
	tags := registry.TagNames()

	// Bundle dynamic import targets (e.g. CSS modules) as preloaded modules.
	// These are specifiers found in thin module import() calls that need to
	// be resolvable at runtime.
	var dynamicBundles []string
	allPreloadModules := append([]string(nil), preloadModules...)
	for _, target := range registry.DynamicImportTargets() {
		modPath, err := ResolveModulePath(target, ".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: cannot resolve dynamic import target %s: %v\n", target, err)
			continue
		}
		bundle, err := BundlePreload(modPath, target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "golit: warning: cannot bundle dynamic import target %s: %v\n", target, err)
			continue
		}
		dynamicBundles = append(dynamicBundles, bundle)
		allPreloadModules = append(allPreloadModules, target)
	}

	drained := make([]*Engine, 0, p.size)

	for i := 0; i < p.size; i++ {
		e := <-p.engines
		e.SetPreloadModules(allPreloadModules)
		e.SetRuntimeExternals(registry.RuntimeExternals())
		for _, pb := range preloadBundles {
			_ = e.LoadBundle(pb)
		}
		for _, db := range dynamicBundles {
			_ = e.LoadBundle(db)
		}
		// Load shared runtime once per engine before any components.
		if rt := registry.SharedRuntime(); rt != "" && !e.loaded["@golit/runtime"] {
			if err := e.LoadModule("@golit/runtime", rt); err != nil {
				fmt.Fprintf(os.Stderr, "golit: warning: loading shared runtime: %v\n", err)
			} else {
				e.loaded["@golit/runtime"] = true
			}
		}
		for _, tag := range tags {
			if _, err := e.LoadBundleForTag(tag, registry); err != nil {
				fmt.Fprintf(os.Stderr, "golit: warning: %v\n", err)
			}
		}
		drained = append(drained, e)
	}

	for _, e := range drained {
		p.engines <- e
	}
	return nil
}

// Get checks out an engine from the pool. Blocks if none are available.
func (p *EnginePool) Get() *Engine {
	return <-p.engines
}

// Put returns an engine to the pool.
func (p *EnginePool) Put(e *Engine) {
	p.engines <- e
}

// Size returns the number of engines in the pool.
func (p *EnginePool) Size() int {
	return p.size
}

// Close releases all engines in the pool.
func (p *EnginePool) Close() {
	close(p.engines)
	for e := range p.engines {
		e.Close()
	}
}
