package data_store

import (
	"gitlab.altex.ro/ams/go_ddwcs/awb/net/rest"
	commonDataStore "gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	lib "gitlab.altex.ro/plug/go_atx_lib"
	cache "gitlab.altex.ro/plug/go_atx_lib/go_redis"
	"os"
)

var (
	localConnection = func() *cache.RedisConnection {
		return &cache.RedisConnection{
			Name:     "internal_instance",
			Url:      os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
		}
	}
)

/**
 * Creates a data store list which contains the available order awbs stores
 * Currently we use Redis (cache) and Oms web service
 */
func GetDataStoreList() []commonDataStore.DataStoreReadOnly {
	var stores []commonDataStore.DataStoreReadOnly

	// cache layer, external oms web service
	stores = append(stores, CreateCacheOrderStoreInstance(), CreateOmsExternalStore()) // beware at the stores order

	return stores
}

func CreateOmsExternalStore() commonDataStore.DataStoreReadOnly {
	omsClient, _ := rest.InitOmsClient(lib.GetContainer().GetLogger())
	store := GetExternalStore(omsClient)
	store.SetLogger(lib.GetContainer().GetLogger())
	return store
}

//Creates a general purposes CacheStore
func CreateCacheStoreInstance() *commonDataStore.CacheStore {
	return commonDataStore.GetCacheStore(lib.GetContainer().GetRedisConnector().GetRedisClient(localConnection()))
}

//Creates OrderAwbs CacheStore
func CreateCacheOrderStoreInstance() commonDataStore.DataStore {
	return GetCacheOrderAwbStore(lib.GetContainer().GetRedisConnector().GetRedisClient(localConnection()))
}

//Creates AwbStatusDetails CacheStore
func CreateCacheAwbStatusStoreInstance() *CacheAwbsStatusStore {
	return GetCachedStatusStore(lib.GetContainer().GetRedisConnector().GetRedisClient(localConnection()))
}
