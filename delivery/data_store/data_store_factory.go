package data_store

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	storeCache "gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/cache"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/cms"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/constraints"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/global_config"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/net/rest"
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

	omsConnection = func() *cache.RedisConnection {
		return &cache.RedisConnection{
			Name:     "oms_instance",
			Url:      os.Getenv("OMS_REDIS_URL"),
			Password: os.Getenv("OMS_REDIS_PASSWORD"),
		}
	}
)

//Creates a general purposes CacheStore - local redis instance
func CreateCacheStoreInstance() *data_store.CacheStore {
	return data_store.GetCacheStore(lib.GetContainer().GetRedisConnector().GetRedisClient(localConnection()))
}

// external redis instance - shared with OMS
func CreateOmsCacheStoreInstance() *data_store.CacheStore {
	return data_store.GetCacheStore(lib.GetContainer().GetRedisConnector().GetRedisClient(omsConnection()))
}

func CreateWarehouseNrOfConfirmationsStore() *storeCache.WarehouseNrOfConfirmations {
	cacheInstance := CreateOmsCacheStoreInstance()
	return storeCache.NewWarehouseNrOfConfirmations(cacheInstance)
}

func CreateStocksExternalStore(websiteCode string) *stocks.ExternalStocksStore {
	stockClient, _ := rest.NewStocksClient(websiteCode, lib.GetContainer().GetLogger())
	return stocks.NewExternalStocksStore(stockClient)
}

// Create external warehouse store and configure it
func CreateWarehousesExternalStore(websiteCode string) *stocks.ExternalWarehousesStore {
	stockClient, _ := rest.NewStocksClient(websiteCode, lib.GetContainer().GetLogger())
	store := stocks.NewExternalWarehousesStore(stockClient)
	store.SetLogger(lib.GetContainer().GetLogger())

	return store
}

// Create an instance of the decorated Warehouse Store (cache + external data store)
func CreateDecoratedWarehousesExternalStore(websiteCode string) constraints.ReadWarehouseStoreInterface {
	stockClient := CreateWarehousesExternalStore(websiteCode)
	cacheClient := storeCache.NewCacheWarehouseStore(CreateCacheStoreInstance(), websiteCode)
	cacheClient.SetLogger(lib.GetContainer().GetLogger())
	return stocks.NewCachingWarehouseStoreDecorator(stockClient, cacheClient)
}

// Create external Global Config store and configure it
func CreateGlobalConfigExternalStore() *global_config.ExternalConfigStore {
	gcClient, _ := rest.NewGlobalConfigClient(lib.GetContainer().GetLogger())
	store := global_config.NewExternalConfigStore(gcClient)
	store.SetLogger(lib.GetContainer().GetLogger())

	return store
}

// Create an instance of the decorated Global Config Store (cache + external data store)
func CreateDecoratedGlobalConfigExternalStore() constraints.ReadGlobalConfigStoreInterface {
	gcStore := CreateGlobalConfigExternalStore()
	cacheStore := storeCache.NewCacheGlobalConfigStore(CreateCacheStoreInstance())
	cacheStore.SetLogger(lib.GetContainer().GetLogger())

	return global_config.NewCachingGlobalConfigStoreDecorator(gcStore, cacheStore)
}

//Creates CmsCacheStore
func CreateCacheCmsStoreInstance() *cms.CacheCmsStore {
	return cms.GetCacheCmsStore(lib.GetContainer().GetRedisConnector().GetRedisClient(localConnection()))
}

// Create external Cms external store and configure it
func CreateCmsExternalStore(websiteCode string) *cms.ExternalCmsStore {
	logger := lib.GetContainer().GetLogger()
	cmsClient, _ := rest.NewCMSClient(websiteCode, logger)
	storeCmsStore := cms.NewExternalStocksStore(cmsClient)

	return storeCmsStore
}
