package jsengine

import "sync"

// SharedRenderCache is a thread-safe render result cache shared across
// engines in a pool. It acts as an L2 cache — engines check their local
// cache first (no locking), then fall back to the shared cache.
type SharedRenderCache struct {
	mu    sync.RWMutex
	cache map[string]renderCacheEntry
}

// NewSharedRenderCache creates a new shared render cache.
func NewSharedRenderCache() *SharedRenderCache {
	return &SharedRenderCache{
		cache: make(map[string]renderCacheEntry),
	}
}

func (c *SharedRenderCache) get(key string) (renderCacheEntry, bool) {
	c.mu.RLock()
	entry, ok := c.cache[key]
	c.mu.RUnlock()
	return entry, ok
}

func (c *SharedRenderCache) set(key string, entry renderCacheEntry) {
	c.mu.Lock()
	c.cache[key] = entry
	c.mu.Unlock()
}

// Len returns the number of entries in the shared cache.
func (c *SharedRenderCache) Len() int {
	c.mu.RLock()
	n := len(c.cache)
	c.mu.RUnlock()
	return n
}
