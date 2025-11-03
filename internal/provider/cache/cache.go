package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheProvider interface {
	SetObject(ctx context.Context, key string, data []byte, ttl time.Duration) error

	GetObject(ctx context.Context, key string) ([]byte, error)

	Del(ctx context.Context, key string) error

	SetString(ctx context.Context, key, str string, ttl time.Duration) error

	GetString(ctx context.Context, key string) (string, error)
}

type cacheProviderImpl struct {
	rdb *redis.Client
}

func NewCacheProvider(rdb *redis.Client) CacheProvider {
	return &cacheProviderImpl{rdb}
}

func (c *cacheProviderImpl) SetString(ctx context.Context, key, str string, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, str, ttl).Err()
}

func (c *cacheProviderImpl) GetString(ctx context.Context, key string) (string, error) {
	str, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return str, nil
}

func (c *cacheProviderImpl) SetObject(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

func (c *cacheProviderImpl) GetObject(ctx context.Context, key string) ([]byte, error) {
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func (c *cacheProviderImpl) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
