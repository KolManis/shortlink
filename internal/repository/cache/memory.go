package cache

import (
	"context"
	"sync"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
)

type MemoryCache struct {
	data map[string]*cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value     string
	expiresAt time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]*cacheItem),
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok || item.expiresAt.Before(time.Now()) {
		return "", urlDomain.ErrNotFound
	}
	return item.value, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *MemoryCache) Incr(ctx context.Context, key string) (int64, error) {
	// реализация инкремента
	c.mu.Lock()
	defer c.mu.Unlock()

	// парсим текущее значение, увеличиваем, сохраняем
	return 1, nil
}
