package jsengine

import (
	"sync"
	"sync/atomic"
)

const defaultSharedCacheMax = 4096

// sharedRenderCache is a thread-safe render result cache shared across
// engines in a pool. It acts as an L2 cache — engines check their local
// cache first (no locking), then fall back to the shared cache.
// Capped at maxEntries to prevent unbounded memory growth when component
// attribute cardinality is high.
type sharedRenderCache struct {
	mu         sync.RWMutex
	cache      map[string]renderCacheEntry
	count      atomic.Int32
	maxEntries int
}

func newSharedRenderCache() *sharedRenderCache {
	return &sharedRenderCache{
		cache:      make(map[string]renderCacheEntry),
		maxEntries: defaultSharedCacheMax,
	}
}

func (c *sharedRenderCache) get(key string) (renderCacheEntry, bool) {
	if c.count.Load() == 0 {
		return renderCacheEntry{}, false
	}
	c.mu.RLock()
	entry, ok := c.cache[key]
	c.mu.RUnlock()
	return entry, ok
}

func (c *sharedRenderCache) set(key string, entry renderCacheEntry) {
	if int(c.count.Load()) >= c.maxEntries {
		return
	}
	c.mu.Lock()
	if _, exists := c.cache[key]; !exists {
		if len(c.cache) >= c.maxEntries {
			c.mu.Unlock()
			return
		}
		c.count.Add(1)
	}
	c.cache[key] = entry
	c.mu.Unlock()
}

func (c *sharedRenderCache) len() int {
	return int(c.count.Load())
}
