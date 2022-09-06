package cache

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/global_config"
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/constraints"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
	"encoding/json"
	"time"
)

const GC_CACHE_KEY_PREFIX = "global_config:"
const GC_CACHE_TTL = 30 // cache for 30 minutes

type ReadWriteGlobalConfigStoreInterface interface {
	constraints.ReadGlobalConfigStoreInterface
	Save(key string, response global_config.ConfigGet)
}

type CacheGlobalConfigStore struct {
	cacheStore  *data_store.CacheStore
	websiteCode string
	logger.LoggerAware
}

func NewCacheGlobalConfigStore(cacheStore *data_store.CacheStore) *CacheGlobalConfigStore {
	cache := &CacheGlobalConfigStore{
		cacheStore: cacheStore,
	}

	return cache
}

// Get Cached GlobalConfig response
// @param filters map[string]string
// @return global_config.ConfigGet
// @desc - returns GlobalConfig cached response
//
func (cg *CacheGlobalConfigStore) GetConfig(filters map[string]string) global_config.ConfigGet {
	var response global_config.ConfigGet
	key := GC_CACHE_KEY_PREFIX + util.Implode(filters)
	result := cg.cacheStore.Get(key)
	if result == "" {
		return response
	}

	if err := json.Unmarshal([]byte(result), &response); err != nil {
		return response
	}

	return response
}

// Get Cached GlobalConfig response
// @param key string
// @param response global_config.ConfigGet
// @desc - store GlobalConfig response in cache for GC_CACHE_TTL min
//
func (cg *CacheGlobalConfigStore) Save(key string, response global_config.ConfigGet) {
	if cg.cacheStore == nil {
		return
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		cg.GetLogger().Println(err.Error())
		return
	}

	cg.cacheStore.Save(data_store.Param{
		Key:   GC_CACHE_KEY_PREFIX + key,
		Value: string(encoded),
		Ttl:   time.Minute * GC_CACHE_TTL,
	})
}
