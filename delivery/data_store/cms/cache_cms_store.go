package cms

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	commonDataStore "gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/cms"
	"time"
)

type CacheCmsStore struct {
	commonDataStore.CacheStore
	client *redis.Client
}

func GetCacheCmsStore(client *redis.Client) *CacheCmsStore {
	cache := &CacheCmsStore{
		client: client,
	}

	// add parent struct dependencies
	cache.CacheStore.SetClient(client)

	return cache
}

/**
 * Get
 * @param key string
 * @return []cms.Store
 * @desc - get the cached stores
 */
func (cacheStore *CacheCmsStore) Get(key string) []cms.Store {
	result := cacheStore.CacheStore.Get(key)

	if result == "" {
		return nil
	}

	var stores []cms.Store
	err := json.Unmarshal([]byte(result), &stores)
	if err != nil {
		return nil
	}

	return stores
}

/**
 * Save
 * @param key string
 * @param info []cms.Store
 * @return void
 * @desc - saves stores for 15 minutes
 */
func (cacheStore *CacheCmsStore) Save(key string, info []cms.Store, ttl time.Duration) {
	if info == nil {
		return
	}

	toCache, err := json.Marshal(info)
	if err != nil {
		return
	}

	cacheStore.CacheStore.Save(commonDataStore.Param{Key: key, Value: toCache, Ttl: ttl})
}
