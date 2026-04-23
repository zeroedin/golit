package jsengine

import (
	"fmt"
	"os"
)

// EnginePool manages a fixed set of QJS Engine instances for concurrent use.
// Each engine is an isolated WASM module instance (via wazero) and can safely
// run in its own goroutine.
type EnginePool struct {
	engines     chan *Engine
	size        int
	sharedCache *SharedRenderCache
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
	if size > 1 {
		p.sharedCache = NewSharedRenderCache()
	}
	for i := 0; i < size; i++ {
		e, err := NewEngine()
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("creating engine %d/%d: %w", i+1, size, err)
		}
		e.sharedCache = p.sharedCache
		p.engines <- e
	}
	return p, nil
}

// PreloadAll drains every engine from the pool, configures preload modules,
// loads any raw preload bundles, bundles dynamic import targets as preloads,
// loads the shared runtime (if present), then loads all registry component
// bundles/modules.
// A dedicated compiler engine compiles modules to bytecode and stores them
// in the registry; subsequent engines load from bytecode for faster initialization.
// After this call the registry must be treated as read-only.
func (p *EnginePool) PreloadAll(registry *Registry, preloadModules []string, preloadBundles ...string) error {
	tags := registry.TagNames()
	runtimeExternals := registry.RuntimeExternals()
	sharedRuntime := registry.SharedRuntime()
	dynamicModules := registry.DynamicModules()

	drained := make([]*Engine, 0, p.size)

	for i := 0; i < p.size; i++ {
		e := <-p.engines
		useBytecode := i > 0 && registry.HasBytecode()

		e.SetPreloadModules(preloadModules)
		e.SetRuntimeExternals(runtimeExternals)
		for _, pb := range preloadBundles {
			_ = e.LoadBundle(pb)
		}
		if sharedRuntime != "" && !e.loaded["@golit/runtime"] {
			if useBytecode {
				if bc := registry.LookupBytecode("@golit/runtime"); bc != nil {
					if err := e.LoadModuleBytecode("@golit/runtime", bc); err == nil {
						e.loaded["@golit/runtime"] = true
					} else {
						fmt.Fprintf(os.Stderr, "golit: bytecode load failed for runtime, falling back to source: %v\n", err)
						useBytecode = false
					}
				}
			}
			if !e.loaded["@golit/runtime"] {
				if err := e.LoadModule("@golit/runtime", sharedRuntime); err != nil {
					fmt.Fprintf(os.Stderr, "golit: warning: loading shared runtime: %v\n", err)
				} else {
					e.loaded["@golit/runtime"] = true
				}
			}
		}
		for specifier, source := range dynamicModules {
			if !e.loaded[specifier] {
				if useBytecode {
					if bc := registry.LookupBytecode(specifier); bc != nil {
						if err := e.LoadModuleBytecode(specifier, bc); err == nil {
							e.loaded[specifier] = true
							continue
						}
					}
				}
				if err := e.LoadModule(specifier, source); err != nil {
					fmt.Fprintf(os.Stderr, "golit: warning: loading dynamic module %s: %v\n", specifier, err)
				} else {
					e.loaded[specifier] = true
				}
			}
		}
		for _, tag := range tags {
			if e.loaded[tag] {
				continue
			}
			if useBytecode {
				if bc := registry.LookupBytecode(tag + ".js"); bc != nil {
					if err := e.EvalModuleBytecode(tag+".js", bc); err == nil {
						e.loaded[tag] = true
						continue
					}
				}
			}
			if _, err := e.LoadBundleForTag(tag, registry); err != nil {
				fmt.Fprintf(os.Stderr, "golit: warning: %v\n", err)
			}
		}

		// First iteration: use a dedicated compiler engine to build bytecode for subsequent engines
		if i == 0 && p.size > 1 && !registry.HasBytecode() {
			p.compileBytecodeCache(e, registry, sharedRuntime, dynamicModules, tags)
		}

		drained = append(drained, e)
	}

	for _, e := range drained {
		p.engines <- e
	}
	return nil
}

// compileBytecodeCache uses a dedicated engine to compile all modules
// to bytecode and store them in the registry for reuse.
func (p *EnginePool) compileBytecodeCache(e *Engine, registry *Registry, sharedRuntime string, dynamicModules map[string]string, tags []string) {
	compiler, err := NewEngine()
	if err != nil {
		return
	}
	defer compiler.Close()

	compiler.SetPreloadModules(e.preloadModules)
	compiler.SetRuntimeExternals(e.runtimeExternals)

	if sharedRuntime != "" {
		if bc, err := compiler.CompileModule("@golit/runtime", sharedRuntime); err == nil {
			registry.SetBytecode("@golit/runtime", bc)
		}
		if err := compiler.LoadModule("@golit/runtime", sharedRuntime); err == nil {
			compiler.loaded["@golit/runtime"] = true
		} else {
			fmt.Fprintf(os.Stderr, "golit: bytecode compiler: failed to load runtime: %v\n", err)
		}
	}

	for specifier, source := range dynamicModules {
		if bc, err := compiler.CompileModule(specifier, source); err == nil {
			registry.SetBytecode(specifier, bc)
		}
		if err := compiler.LoadModule(specifier, source); err == nil {
			compiler.loaded[specifier] = true
		} else {
			fmt.Fprintf(os.Stderr, "golit: bytecode compiler: failed to load %s: %v\n", specifier, err)
		}
	}

	for _, tag := range tags {
		source := registry.Lookup(tag)
		if source == "" {
			continue
		}
		if bc, err := compiler.CompileModule(tag+".js", source); err == nil {
			registry.SetBytecode(tag+".js", bc)
		}
	}
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

// Available returns the number of idle engines currently in the pool.
func (p *EnginePool) Available() int {
	return len(p.engines)
}

// Close releases all engines in the pool.
func (p *EnginePool) Close() {
	close(p.engines)
	for e := range p.engines {
		e.Close()
	}
}
