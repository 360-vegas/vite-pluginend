package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 缓存接口
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Close() error
}

// RedisCache Redis 缓存实现
type RedisCache struct {
	client Cache
}

// NewRedisCache 创建新的 Redis 缓存实例
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: &redisClient{client: client}}, nil
}

// NewMemoryCache 创建新的内存缓存实例
func NewMemoryCache() *RedisCache {
	return &RedisCache{client: &MemoryCache{
		data: make(map[string][]byte),
	}}
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration)
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	return c.client.Get(ctx, key, value)
}

// Delete 删除缓存
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Delete(ctx, key)
}

// Close 关闭连接
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// redisClient Redis 客户端实现
type redisClient struct {
	client *redis.Client
}

// Set 设置缓存
func (c *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *redisClient) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

// Delete 删除缓存
func (c *redisClient) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Close 关闭连接
func (c *redisClient) Close() error {
	return c.client.Close()
} 