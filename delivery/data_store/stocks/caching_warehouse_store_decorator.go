package stocks

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/cache"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store/constraints"
)

type CachingWarehouseStoreDecorator struct {
	dataStore  constraints.ReadWarehouseStoreInterface
	cacheStore cache.ReadWriteWarehouseStoreInterface
}

func NewCachingWarehouseStoreDecorator(dataStore constraints.ReadWarehouseStoreInterface, cacheStore cache.ReadWriteWarehouseStoreInterface) constraints.ReadWarehouseStoreInterface {
	return &CachingWarehouseStoreDecorator{
		dataStore:  dataStore,
		cacheStore: cacheStore,
	}
}

/**
 * GetWarehouses
 * @param map[string]string
 * return []stocks.Warehouse
 */
func (decorator *CachingWarehouseStoreDecorator) GetWarehouses(filters map[string]string) []stocks.Warehouse {
	cached := decorator.cacheStore.GetWarehouses(filters)
	if len(cached) > 0 {
		return cached
	}

	whs := decorator.dataStore.GetWarehouses(filters)
	decorator.cacheStore.Save(decorator.prepareCacheKey(filters), whs)

	return whs
}

/**
 * prepareCacheKey - private
 * @param map[string]string
 * return string
 * @desc - exclude some keys and implode the useful ones
 */
func (decorator *CachingWarehouseStoreDecorator) prepareCacheKey(filters map[string]string) string {
	delete(filters, "page_no")
	delete(filters, "items_per_page")
	return util.Implode(filters)
}
