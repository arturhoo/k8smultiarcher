package main

import (
	"context"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	cacheSizeDefault = 100000
	redisAddrDefault = "localhost:6379"
)

type Cache interface {
	Get(key string) (bool, bool)
	Set(key string, value bool)
}

type InMemoryCache struct {
	cache gcache.Cache
}

func NewInMemoryCache(cacheSize int) *InMemoryCache {
	gc := gcache.New(cacheSize).ARC().Build()
	return &InMemoryCache{gc}
}

func (c InMemoryCache) Get(key string) (bool, bool) {
	val, err := c.cache.Get(key)
	if err != nil {
		return false, false
	}
	boolVal, ok := val.(bool)
	if !ok {
		log.Fatal().Msg("found non boolean cache value")
	}
	return boolVal, true
}

func (c *InMemoryCache) Set(key string, value bool) {
	err := c.cache.Set(key, value)
	if err != nil {
		log.Printf("got error trying to set key on InMemoryCache: %s", err)
	}
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisAddr string) *RedisCache {
	return &RedisCache{
		redis.NewClient(&redis.Options{
			Addr: redisAddr,
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
