package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache interface {
	Set(key string, value interface{})
	Get(key string) string
	Del(key string)
}

type redisCache struct {
	host     string
	db       int
	password string
	expires  time.Duration
}

func NewRedisCache(host string, db int, password string, expires time.Duration) RedisCache {
	return &redisCache{
		host:     host,
		db:       db,
		password: password,
		expires:  expires,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: cache.password,
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value interface{}) {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	if err = client.Set(key, json, cache.expires*time.Minute).Err(); err != nil {
		panic(err)
	}
}

func (cache *redisCache) Get(key string) (email string) {
	client := cache.getClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return ""
	}

	if err := json.Unmarshal([]byte(val), &email); err != nil {
		panic(err)
	}
	return
}

func (cache *redisCache) Del(key string) {
	client := cache.getClient()
	defer client.Close()

	client.Del(key)
}
