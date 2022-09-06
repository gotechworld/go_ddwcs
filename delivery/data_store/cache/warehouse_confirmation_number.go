package cache

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"strconv"
	"time"
)

type WarehouseNrOfConfirmations struct {
	cache *data_store.CacheStore
}

// WarehouseNrOfConfirmations - constructor
func NewWarehouseNrOfConfirmations(cache *data_store.CacheStore) *WarehouseNrOfConfirmations {
	return &WarehouseNrOfConfirmations{cache: cache}
}

// Get today number of confirmations for "warehouses"
func (whCon *WarehouseNrOfConfirmations) GetTodayNrOfConfirmations(warehouses *map[int]stocks.Warehouse) map[int]int {
	nrOfConfirmations := map[int]int{}
	for _, warehouse := range *warehouses {
		todayNrOfConfirmations := whCon.cache.Get(whCon.getRedisKey(warehouse.Id))
		nr, err := strconv.Atoi(todayNrOfConfirmations)
		if nil != err {
			nr = 0
		}

		nrOfConfirmations[warehouse.Id] = nr
	}

	return nrOfConfirmations
}

// Get Redis key composed as in OMS confirm process (value at this key is incremented by OMS confirm process)
func (whCon *WarehouseNrOfConfirmations) getRedisKey(warehouseId int) string {
	return "Order_Confirm_noOfConfirmations_" + time.Now().Format("2006-01-02") + "_" + strconv.Itoa(warehouseId)
}
