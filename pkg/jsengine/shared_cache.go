package jsengine

import "sync"

// sharedRenderCache is a thread-safe render result cache shared across
// engines in a pool. It acts as an L2 cache — engines check their local
// cache first (no locking), then fall back to the shared cache.
type sharedRenderCache struct {
	mu    sync.RWMutex
	cache map[string]renderCacheEntry
}

func newSharedRenderCache() *sharedRenderCache {
	return &sharedRenderCache{
		cache: make(map[string]renderCacheEntry),
	}
}

func (c *sharedRenderCache) get(key string) (renderCacheEntry, bool) {
	c.mu.RLock()
	entry, ok := c.cache[key]
	c.mu.RUnlock()
	return entry, ok
}

func (c *sharedRenderCache) set(key string, entry renderCacheEntry) {
	c.mu.Lock()
	c.cache[key] = entry
	c.mu.Unlock()
}

func (c *sharedRenderCache) len() int {
	c.mu.RLock()
	n := len(c.cache)
	c.mu.RUnlock()
	return n
}
