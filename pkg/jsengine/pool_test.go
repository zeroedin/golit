package jsengine

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestEnginePool_GetPut(t *testing.T) {
	pool, err := NewEnginePool(2)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if pool.Size() != 2 {
		t.Fatalf("Size = %d, want 2", pool.Size())
	}

	e1 := pool.Get()
	e2 := pool.Get()

	// Both engines should be usable independently.
	if e1 == e2 {
		t.Error("expected distinct engine instances")
	}

	pool.Put(e1)
	pool.Put(e2)

	// Should be able to get them again.
	e3 := pool.Get()
	pool.Put(e3)
}

func TestEnginePool_PreloadAndRender(t *testing.T) {
	bundle := bundleMyGreeting(t)

	registry := NewRegistry()
	tagName, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	registry.Register(tagName, bundle)

	pool, err := NewEnginePool(2)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if err := pool.PreloadAll(registry, nil); err != nil {
		t.Fatalf("PreloadAll: %v", err)
	}

	// Both engines should be able to render after preloading.
	for i := 0; i < 2; i++ {
		e := pool.Get()
		result, err := e.RenderElement("my-greeting", map[string]string{"name": "Pool"})
		pool.Put(e)

		if err != nil {
			t.Fatalf("render %d: %v", i, err)
		}
		if !strings.Contains(result.HTML, "Pool") {
			t.Errorf("render %d: missing 'Pool' in output", i)
		}
	}
}

func TestEnginePool_Available(t *testing.T) {
	pool, err := NewEnginePool(3)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if got := pool.Available(); got != 3 {
		t.Fatalf("Available = %d, want 3 (all idle)", got)
	}

	e1 := pool.Get()
	if got := pool.Available(); got != 2 {
		t.Fatalf("Available after 1 Get = %d, want 2", got)
	}

	e2 := pool.Get()
	if got := pool.Available(); got != 1 {
		t.Fatalf("Available after 2 Gets = %d, want 1", got)
	}

	pool.Put(e1)
	if got := pool.Available(); got != 2 {
		t.Fatalf("Available after Put = %d, want 2", got)
	}

	pool.Put(e2)
	if got := pool.Available(); got != 3 {
		t.Fatalf("Available after all returned = %d, want 3", got)
	}
}

func TestEnginePool_BytecodePrecompilation(t *testing.T) {
	bundle := bundleMyGreeting(t)

	registry := NewRegistry()
	tagName, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	registry.Register(tagName, bundle)

	pool, err := NewEnginePool(2)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if registry.HasBytecode() {
		t.Fatal("registry should not have bytecode before PreloadAll")
	}

	if err := pool.PreloadAll(registry, nil); err != nil {
		t.Fatalf("PreloadAll: %v", err)
	}

	if !registry.HasBytecode() {
		t.Fatal("registry should have bytecode after PreloadAll with pool size > 1")
	}

	e := pool.Get()
	result, err := e.RenderElement("my-greeting", map[string]string{"name": "Bytecode"})
	pool.Put(e)
	if err != nil {
		t.Fatalf("render with bytecode-loaded engine: %v", err)
	}
	if !strings.Contains(result.HTML, "Bytecode") {
		t.Errorf("missing 'Bytecode' in output: %s", result.HTML)
	}
}

func TestEnginePool_SharedCache_NotCreatedForSingleEngine(t *testing.T) {
	pool, err := NewEnginePool(1)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if pool.sharedCache != nil {
		t.Error("pool of size 1 should not have a shared cache")
	}

	e := pool.Get()
	if e.sharedCache != nil {
		t.Error("engine in pool of size 1 should not have a shared cache")
	}
	pool.Put(e)
}

func TestEnginePool_SharedCache_CreatedForMultipleEngines(t *testing.T) {
	pool, err := NewEnginePool(2)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if pool.sharedCache == nil {
		t.Fatal("pool of size 2 should have a shared cache")
	}

	e1 := pool.Get()
	e2 := pool.Get()

	if e1.sharedCache == nil || e2.sharedCache == nil {
		t.Error("both engines should have a shared cache reference")
	}
	if e1.sharedCache != e2.sharedCache {
		t.Error("both engines should share the same cache instance")
	}

	pool.Put(e1)
	pool.Put(e2)
}

func TestEnginePool_SharedCache_CrossEngineHit(t *testing.T) {
	bundle := bundleMyGreeting(t)

	registry := NewRegistry()
	tagName, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	registry.Register(tagName, bundle)

	pool, err := NewEnginePool(2)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if err := pool.PreloadAll(registry, nil); err != nil {
		t.Fatalf("PreloadAll: %v", err)
	}

	attrs := map[string]string{"name": "Shared"}
	reqs := []BatchRequest{{ID: 1, TagName: "my-greeting", Attrs: attrs}}

	// Engine A renders via RenderBatch and populates the shared cache.
	eA := pool.Get()
	resultsA, err := eA.RenderBatch(reqs)
	if err != nil {
		t.Fatalf("engine A render: %v", err)
	}
	pool.Put(eA)

	if pool.sharedCache.len() == 0 {
		t.Fatal("shared cache should have entries after engine A render")
	}

	// Engine B should get the result from shared cache (L2 hit).
	eB := pool.Get()
	resultsB, err := eB.RenderBatch(reqs)
	if err != nil {
		t.Fatalf("engine B render: %v", err)
	}
	pool.Put(eB)

	if resultsA[0].HTML != resultsB[0].HTML {
		t.Errorf("cross-engine results differ:\n  A: %s\n  B: %s", resultsA[0].HTML, resultsB[0].HTML)
	}
	if resultsA[0].CSS != resultsB[0].CSS {
		t.Errorf("cross-engine CSS differs:\n  A: %s\n  B: %s", resultsA[0].CSS, resultsB[0].CSS)
	}
}

func TestEnginePool_SharedCache_PromotesToL1(t *testing.T) {
	bundle := bundleMyGreeting(t)

	registry := NewRegistry()
	tagName, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	registry.Register(tagName, bundle)

	pool, err := NewEnginePool(2)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if err := pool.PreloadAll(registry, nil); err != nil {
		t.Fatalf("PreloadAll: %v", err)
	}

	attrs := map[string]string{"name": "Promote"}
	reqs := []BatchRequest{{ID: 1, TagName: "my-greeting", Attrs: attrs}}

	// Engine A renders via RenderBatch — populates L1 + L2.
	eA := pool.Get()
	_, err = eA.RenderBatch(reqs)
	if err != nil {
		t.Fatalf("engine A render: %v", err)
	}
	pool.Put(eA)

	// Engine B renders — L2 hit, should promote to B's L1.
	eB := pool.Get()
	_, err = eB.RenderBatch(reqs)
	if err != nil {
		t.Fatalf("engine B first render: %v", err)
	}

	key := renderCacheKey("my-greeting", attrs)
	if _, ok := eB.renderCache[key]; !ok {
		t.Error("shared cache hit should be promoted to engine B's local cache")
	}

	pool.Put(eB)
}

func TestEnginePool_ConcurrentRender(t *testing.T) {
	bundle := bundleMyGreeting(t)

	registry := NewRegistry()
	tagName, err := DiscoverTagName(bundle)
	if err != nil {
		t.Fatalf("DiscoverTagName: %v", err)
	}
	registry.Register(tagName, bundle)

	pool, err := NewEnginePool(3)
	if err != nil {
		t.Fatalf("NewEnginePool: %v", err)
	}
	defer pool.Close()

	if err := pool.PreloadAll(registry, nil); err != nil {
		t.Fatalf("PreloadAll: %v", err)
	}

	names := []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank"}
	var wg sync.WaitGroup
	errors := make(chan error, len(names))

	for _, name := range names {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			e := pool.Get()
			defer pool.Put(e)

			result, err := e.RenderElement("my-greeting", map[string]string{"name": n})
			if err != nil {
				errors <- err
				return
			}
			if !strings.Contains(result.HTML, n) {
				errors <- fmt.Errorf("missing %q in output: %s", n, result.HTML)
			}
		}(name)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}
