package jsengine

import "fmt"

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

// PreloadAll loads every bundle from the registry into every engine in the
// pool, and configures preload modules. After this call the registry must
// be treated as read-only.
func (p *EnginePool) PreloadAll(registry *Registry, preloadModules []string) error {
	tags := registry.TagNames()
	drained := make([]*Engine, 0, p.size)

	// Drain all engines so we can configure each one.
	for i := 0; i < p.size; i++ {
		e := <-p.engines
		e.SetPreloadModules(preloadModules)
		for _, tag := range tags {
			e.LoadBundleForTag(tag, registry)
		}
		drained = append(drained, e)
	}

	// Return them.
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
