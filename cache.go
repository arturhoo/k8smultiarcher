package main

import (
	"context"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
)

const cacheSize = 100000

type Cache interface {
	Get(key string) (bool, bool)
	Set(key string, value bool)
}

type InMemoryCache struct {
	cache gcache.Cache
}

func NewInMemoryCache() *InMemoryCache {
	gc := gcache.New(cacheSize).ARC().Build()
	return &InMemoryCache{gc}
}

func (c InMemoryCache) Get(key string) (bool, bool) {
	val, err := c.cache.Get(key)
	if err != nil {
		return false, false
	}
	return val.(bool), true
}

func (c *InMemoryCache) Set(key string, value bool) {
	c.cache.Set(key, value)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() RedisCache {
	return RedisCache{
		redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func (c RedisCache) Get(key string) (bool, bool) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Bool()
	if err != nil {
		return false, false
	}
	return val, true
}

func (c *RedisCache) Set(key string, value bool) {
	ctx := context.Background()
	c.client.Set(ctx, key, value, 0)
}
