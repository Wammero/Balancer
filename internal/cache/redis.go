package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	ClientCache ClientCache
	ProxyCache  ProxyCache
	client      *redis.Client
	ctx         context.Context
}

func NewRedisCache(host, port string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	return &RedisCache{
		ClientCache: NewClientCache(rdb),
		ProxyCache:  NewProxyCache(rdb),
		client:      rdb,
		ctx:         context.Background(),
	}
}

func (r *RedisCache) Set(key, value string, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(key string) (string, bool, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
