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
