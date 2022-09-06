package cache

import (
	"encoding/json"
	"gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/constraints"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
	"time"
)

const WAREHOUSE_KEY_PREFIX = "warehouse_active_"

type ReadWriteWarehouseStoreInterface interface {
	constraints.ReadWarehouseStoreInterface
	Save(string, []stocks.Warehouse)
}

type CacheWarehouseStore struct {
	client      *data_store.CacheStore
	websiteCode string
	logger.LoggerAware
}

func NewCacheWarehouseStore(client *data_store.CacheStore, websiteCode string) *CacheWarehouseStore {
	cache := &CacheWarehouseStore{
		client:      client,
		websiteCode: websiteCode,
	}

	return cache
}

/**
 * GetWarehouses
 * @param filters map[string]string
 * @return []stocks.Warehouse
 * @desc - returns a list of cached warehouses
 */
func (cache *CacheWarehouseStore) GetWarehouses(filters map[string]string) []stocks.Warehouse {
	var whs []stocks.Warehouse
	key := WAREHOUSE_KEY_PREFIX + cache.websiteCode + util.Implode(filters)
	result := cache.client.Get(key)
	if result == "" {
		return whs
	}

	if err := json.Unmarshal([]byte(result), &whs); err != nil {
		return whs
	}

	return whs
}

/**
 * Save
 * @param whs []stocks.Warehouse
 * @desc - store in cache for 30 minutes
 */
func (cache *CacheWarehouseStore) Save(key string, whs []stocks.Warehouse) {
	if cache.client == nil || len(whs) == 0 {
		return
	}

	encoded, err := json.Marshal(whs)
	if err != nil {
		cache.GetLogger().Println(err.Error())
		return
	}

	cache.client.Save(data_store.Param{
		Key:   WAREHOUSE_KEY_PREFIX + cache.websiteCode + key,
		Value: string(encoded),
		Ttl:   time.Hour / 2, // cache for 30 minutes
	})
}
