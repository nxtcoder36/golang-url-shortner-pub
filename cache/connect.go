package cache

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	Client *redis.Client
	TTL    time.Duration // time to live
}

type RedisCacheInterface interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

func RedisCacheImpl() (RedisCacheInterface, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
	})
	t, err := strconv.Atoi(os.Getenv("REDIS_TTL"))
	if err != nil {
		return nil, err
	}
	return &RedisCache{
		Client: rdb,
		TTL:    time.Duration(t) * time.Second,
	}, nil
}

func (r *RedisCache) Set(key string, value string) error {
	return r.Client.Set(strings.Trim(key, "/"), value, r.TTL).Err()
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.Client.Get(strings.Trim(key, "/")).Result()
}
