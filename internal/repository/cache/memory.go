package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
)

type MemoryCache struct {
	data   map[string]*cacheItem
	mu     sync.RWMutex
	logger *slog.Logger
}

type cacheItem struct {
	value     string
	expiresAt time.Time
}

func NewMemoryCache(logger *slog.Logger) *MemoryCache {
	return &MemoryCache{
		data:   make(map[string]*cacheItem),
		logger: logger,
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	c.logger.Debug("Memory GET", "key", key)

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.data[key]
	if !ok || item.expiresAt.Before(time.Now()) {
		c.logger.Debug("Memory GET miss or Memory GET expired", "key", key)
		return "", urlDomain.ErrNotFound
	}

	c.logger.Debug("Memory GET hit", "key", key, "value", item.value)
	return item.value, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	c.logger.Debug("Memory SET", "key", key, "ttl", ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.logger.Debug("Memory DELETE", "key", key)

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *MemoryCache) Incr(ctx context.Context, key string) (int64, error) {
	c.logger.Debug("Memory INCR", "key", key)

	// реализация инкремента
	c.mu.Lock()
	defer c.mu.Unlock()

	// парсим текущее значение, увеличиваем, сохраняем
	return 1, nil
}
