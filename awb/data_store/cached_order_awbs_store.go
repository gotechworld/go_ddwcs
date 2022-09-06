package data_store

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"gitlab.altex.ro/ams/go_ddwcs/awb/dto"
	"gitlab.altex.ro/ams/go_ddwcs/awb/util"
	commonDataStore "gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	"time"
)

type CacheOrderAwbsStore struct {
	commonDataStore.CacheStore
	client *redis.Client
}

func GetCacheOrderAwbStore(client *redis.Client) *CacheOrderAwbsStore {
	cache := &CacheOrderAwbsStore{
		client: client,
	}

	// add parent struct dependencies
	cache.CacheStore.SetClient(client)

	return cache
}

/**
 * GetList
 * @param params map[string]string
 * @return []dto.AwbInfo
 * @desc - get the cached awbs list
 * the params map will contain the following keys: increment_id, customer_id, customer_email
 */
func (cacheStore *CacheOrderAwbsStore) GetList(params map[string]string) []dto.AwbInfo {
	var emptyList []dto.AwbInfo
	if len(params) == 0 {
		return emptyList
	}

	// extract the key from the provided params
	key := util.CreateKeyFromOrderParams(params)
	val := cacheStore.CacheStore.Get(key)
	if val == "" {
		return emptyList
	}

	list := dto.OrderAwbList{}
	err := json.Unmarshal([]byte(val), &list)
	if err != nil {
		return emptyList
	}

	return list.Awbs
}

/**
 * Save
 * @param key string
 * @param list []dto.AwbInfo
 * @return void
 * @desc - saves a list of awbs details in cache for one hour
 */
func (cacheStore *CacheOrderAwbsStore) Save(key string, list []dto.AwbInfo) {
	orderAwbInfo := dto.OrderAwbList{
		Awbs: list,
	}

	toCache, err := json.Marshal(orderAwbInfo)
	if err != nil {
		return
	}

	cacheStore.CacheStore.Save(commonDataStore.Param{Key: key, Value: toCache, Ttl: time.Hour / 2}) //half an hour expiration
}
