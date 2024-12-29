package jankdb

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache[T] is a thin wrapper around github.com/patrickmn/go-cache,
// storing entire T objects under a single key (like "all").
type Cache[T any] struct {
	c *gocache.Cache
}

func NewCache[T any](defaultExpiration, cleanupInterval time.Duration) *Cache[T] {
	return &Cache[T]{
		c: gocache.New(defaultExpiration, cleanupInterval),
	}
}

// Set stores a T under a given key.
func (cache *Cache[T]) Set(key string, val T) {
	cache.c.Set(key, val, gocache.DefaultExpiration)
}

// Get retrieves the T stored under the key.
func (cache *Cache[T]) Get(key string) (T, bool) {
	value, found := cache.c.Get(key)
	if !found {
		var zero T
		return zero, false
	}
	tVal, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return tVal, true
}

// Delete removes a key from the cache.
func (cache *Cache[T]) Delete(key string) {
	cache.c.Delete(key)
}
