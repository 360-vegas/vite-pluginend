package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// MemoryCache 内存缓存实现
type MemoryCache struct {
	data  map[string][]byte
	mutex sync.RWMutex
}

// Set 设置缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = data

	// 如果设置了过期时间，启动一个 goroutine 来删除过期的数据
	if expiration > 0 {
		go func() {
			time.Sleep(expiration)
			c.mutex.Lock()
			delete(c.data, key)
			c.mutex.Unlock()
		}()
	}

	return nil
}

// Get 获取缓存
func (c *MemoryCache) Get(ctx context.Context, key string, value interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	data, ok := c.data[key]
	if !ok {
		return errors.New("key not found")
	}

	return json.Unmarshal(data, value)
}

// Delete 删除缓存
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
	return nil
}

// Close 关闭缓存
func (c *MemoryCache) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = make(map[string][]byte)
	return nil
} 