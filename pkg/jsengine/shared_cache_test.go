package jsengine

import (
	"sync"
	"testing"
)

func TestSharedRenderCache_GetSet(t *testing.T) {
	c := NewSharedRenderCache()

	if _, ok := c.get("missing"); ok {
		t.Error("expected miss for absent key")
	}

	c.set("k1", renderCacheEntry{html: "<h1>hi</h1>", css: ".a{}"})
	entry, ok := c.get("k1")
	if !ok {
		t.Fatal("expected hit after set")
	}
	if entry.html != "<h1>hi</h1>" || entry.css != ".a{}" {
		t.Errorf("unexpected entry: %+v", entry)
	}
}

func TestSharedRenderCache_Len(t *testing.T) {
	c := NewSharedRenderCache()
	if c.Len() != 0 {
		t.Fatalf("Len = %d, want 0", c.Len())
	}

	c.set("a", renderCacheEntry{html: "a"})
	c.set("b", renderCacheEntry{html: "b"})
	if c.Len() != 2 {
		t.Fatalf("Len = %d, want 2", c.Len())
	}

	c.set("a", renderCacheEntry{html: "a2"})
	if c.Len() != 2 {
		t.Fatalf("Len after overwrite = %d, want 2", c.Len())
	}
}

func TestSharedRenderCache_ConcurrentAccess(t *testing.T) {
	c := NewSharedRenderCache()
	var wg sync.WaitGroup
	n := 100

	for i := 0; i < n; i++ {
		wg.Add(2)
		key := string(rune('a' + (i % 26)))
		go func(k string, v int) {
			defer wg.Done()
			c.set(k, renderCacheEntry{html: "<p>" + k + "</p>"})
		}(key, i)
		go func(k string) {
			defer wg.Done()
			c.get(k)
		}(key)
	}

	wg.Wait()
	if c.Len() == 0 {
		t.Error("expected some entries after concurrent writes")
	}
}
