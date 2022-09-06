package global_config

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/global_config"
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/cache"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/constraints"
)

type CachingGlobalConfigStoreDecorator struct {
	dataStore  constraints.ReadGlobalConfigStoreInterface
	cacheStore cache.ReadWriteGlobalConfigStoreInterface
}

func NewCachingGlobalConfigStoreDecorator(dataStore constraints.ReadGlobalConfigStoreInterface, cacheStore cache.ReadWriteGlobalConfigStoreInterface) constraints.ReadGlobalConfigStoreInterface {
	return &CachingGlobalConfigStoreDecorator{
		dataStore:  dataStore,
		cacheStore: cacheStore,
	}
}

// GetConfig
// @param map[string]string
// return global_config.ConfigGet
//
func (decorator *CachingGlobalConfigStoreDecorator) GetConfig(params map[string]string) global_config.ConfigGet {
	cached := decorator.cacheStore.GetConfig(params)
	if len(cached.ConfigGetData) > 0 {
		return cached
	}

	response := decorator.dataStore.GetConfig(params)
	decorator.cacheStore.Save(decorator.prepareCacheKey(params), response)

	return response
}

// prepareCacheKey - private
// @param map[string]string
// return string
//
func (decorator *CachingGlobalConfigStoreDecorator) prepareCacheKey(params map[string]string) string {
	return util.Implode(params)
}
