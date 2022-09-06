package data_store

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"gitlab.altex.ro/ams/go_ddwcs/awb/service/adapter"
	commonDataStore "gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	"time"
)

type CacheAwbsStatusStore struct {
	commonDataStore.CacheStore
	client *redis.Client
}

func GetCachedStatusStore(client *redis.Client) *CacheAwbsStatusStore {
	cache := &CacheAwbsStatusStore{
		client: client,
	}

	// add parent struct dependencies
	cache.CacheStore.SetClient(client)

	return cache
}

/**
 * Get
 * @param key string
 * @return adapter.AwbStatusDetails pointer
 * @desc - get the cached awbs status details
 * the params map will contain the following keys: increment_id, customer_id, customer_email
 */
func (cacheStore *CacheAwbsStatusStore) Get(key string) *adapter.AwbStatusDetails {
	result := cacheStore.CacheStore.Get(key)

	if result == "" {
		return nil
	}

	var status *adapter.AwbStatusDetails
	err := json.Unmarshal([]byte(result), &status)
	if err != nil {
		return nil
	}

	return status
}

/**
 * Save
 * @param key string
 * @param info adapter.AwbStatusDetails pointer
 * @return void
 * @desc - saves awbs details in cache for one hour
 */
func (cacheStore *CacheAwbsStatusStore) Save(key string, info *adapter.AwbStatusDetails) {
	if info == nil {
		return
	}

	toCache, err := json.Marshal(info)
	if err != nil {
		return
	}

	cacheStore.CacheStore.Save(commonDataStore.Param{Key: key, Value: toCache, Ttl: time.Hour}) //1h expiration
}
