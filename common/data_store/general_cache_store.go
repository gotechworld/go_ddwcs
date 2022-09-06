package data_store

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type Param struct {
	Key   string
	Value interface{}
	Ttl   time.Duration
}

type CacheStore struct {
	client *redis.Client
}

func GetCacheStore(client *redis.Client) *CacheStore {
	return &CacheStore{client: client}
}

/**
 * SetClient
 * @param client *redis.Client
 * @desc - offers the possibility to set the redis client if the client does not uses the constructor
 */
func (cache *CacheStore) SetClient(client *redis.Client) {
	cache.client = client
}

/**
 * Get
 * @param key string
 * @return string
 * @desc - get the cached key result
 */
func (cache *CacheStore) Get(key string) string {
	if cache.client == nil || key == "" {
		return ""
	}

	val, err := cache.client.Get(key).Result()
	if err != nil {
		return ""
	}

	return val
}

/**
 * Get
 * @param param Param
 * @return void
 * @desc - Saves the param.Value under param.Key as a string value - REDIS SETEX
 */
func (cache *CacheStore) Save(param Param) {
	if cache.client == nil || param.Value == nil || param.Key == "" {
		return
	}

	cache.client.Set(param.Key, param.Value, param.Ttl)
}
