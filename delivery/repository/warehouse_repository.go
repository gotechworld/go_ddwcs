package repository

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store"
	"gitlab.altex.ro/plug/go_atx_lib/util"
	logger "gitlab.altex.ro/plug/go_logger"
)

type WarehouseRepository struct {
	// only active warehouses (status = 1)
	warehouses *map[int]stocks.Warehouse
	// active warehouses that don't have threshold reached (noOfTodayConfirmations < threshold)
	availableWarehouses *map[int]stocks.Warehouse
	// today nr of confirmations todayNrOfConfirmations[warehouseId] = todayNumberOfConfirmations
	todayNrOfConfirmations *map[int]int
	// general logger
	logger logger.Logger
}

// WarehouseRepository - constructor
func NewWarehouseRepository(websiteCode string, logger logger.Logger) *WarehouseRepository {
	repo := &WarehouseRepository{}
	whs := repo.loadWarehousesByWebsite(websiteCode)

	whTodayNrOfConfirmations := data_store.CreateWarehouseNrOfConfirmationsStore()
	todayNrOfConfirmations := whTodayNrOfConfirmations.GetTodayNrOfConfirmations(&whs)

	availableWarehouses := map[int]stocks.Warehouse{}
	for _, warehouse := range whs {
		if todayNrOfConfirmations[warehouse.Id] < warehouse.AutoconfirmThreshold {
			availableWarehouses[warehouse.Id] = warehouse
		} else {
			logger.Printf(
				"[WarehouseRepository] removed %v[%v] due to threshold %v reached today :%v",
				warehouse.Name,
				warehouse.Id,
				warehouse.AutoconfirmThreshold,
				todayNrOfConfirmations[warehouse.Id],
			)
		}
	}

	repo.warehouses = &whs
	repo.availableWarehouses = &availableWarehouses
	repo.todayNrOfConfirmations = &todayNrOfConfirmations
	repo.logger = logger

	return repo
}

// loadWarehousesByWebsite - private
// @param websiteCode string
// @return  map[int]stocks.Warehouse
// @desc - load the warehouses from cache or external stock system
//
func (repo *WarehouseRepository) loadWarehousesByWebsite(websiteCode string) map[int]stocks.Warehouse {
	warehousesFromApi := data_store.CreateDecoratedWarehousesExternalStore(websiteCode).
		GetWarehouses(map[string]string{"status": stocks.WH_IS_ACTIVE})
	warehouses := map[int]stocks.Warehouse{}
	for _, warehouse := range warehousesFromApi {
		warehouses[warehouse.Id] = warehouse
	}

	return warehouses
}

// Get all "active"/"enabled" warehouses (status = 1)
func (repo *WarehouseRepository) GetAll() *map[int]stocks.Warehouse {
	return repo.warehouses
}

// Get "active"/"enabled" warehouses that don't have threshold reached (noOfTodayConfirmations < threshold)
func (repo *WarehouseRepository) GetAvailable() *map[int]stocks.Warehouse {
	return repo.availableWarehouses
}

// Get "active"/"enabled" warehouses today nr of confirmations
// map[warehouseId] = todayNumberOfConfirmations
func (repo *WarehouseRepository) GetTodayNrOfConfirmations() *map[int]int {
	return repo.todayNrOfConfirmations
}

// Get "active"/"enabled" warehouses that are in specified region (Proximity Rule)
func (repo *WarehouseRepository) GetInRegion(region string) *map[int]stocks.Warehouse {
	inRegionWarehouses := map[int]stocks.Warehouse{}
	for _, warehouse := range *repo.availableWarehouses {
		if util.InArray(warehouse.ConfirmRegions, region) {
			inRegionWarehouses[warehouse.Id] = warehouse
		}

	}

	return &inRegionWarehouses
}

// Get "active"/"enabled" warehouses that are "zonal"  and in specified region
func (repo *WarehouseRepository) GetZonal(region string) *map[int]stocks.Warehouse {
	inRegionWarehouses := map[int]stocks.Warehouse{}
	for _, warehouse := range *repo.availableWarehouses {
		if warehouse.IsZonal && util.InArray(warehouse.ConfirmRegions, region) {
			inRegionWarehouses[warehouse.Id] = warehouse
		}

	}

	return &inRegionWarehouses
}

// Get "active"/"enabled" warehouses that are "stores" and in specified region
func (repo *WarehouseRepository) GetStoresInRegion(regionCode string) *map[int]stocks.Warehouse {
	inRegionStores := map[int]stocks.Warehouse{}
	for _, warehouse := range *repo.availableWarehouses {
		if 1 == warehouse.IsStore && util.InArray(warehouse.ConfirmRegions, regionCode) {
			inRegionStores[warehouse.Id] = warehouse
		}
	}

	return &inRegionStores
}

// Get "active"/"enabled" warehouses that are "stores" and in one of specified regions
func (repo *WarehouseRepository) GetStoresInRegions(regionCodes []string) *map[int]stocks.Warehouse {
	inRegionStores := map[int]stocks.Warehouse{}
	for _, warehouse := range *repo.availableWarehouses {
		for _, regionCode := range regionCodes {
			if 1 == warehouse.IsStore && util.InArray(warehouse.ConfirmRegions, regionCode) {
				inRegionStores[warehouse.Id] = warehouse
			}
		}
	}

	return &inRegionStores
}

// Get "active"/"enabled" warehouses that are "central" also known as "primary"
func (repo *WarehouseRepository) GetCentral() *map[int]stocks.Warehouse {
	centralWarehouses := map[int]stocks.Warehouse{}
	for _, warehouse := range *repo.availableWarehouses {
		if stocks.WH_IS_PRIMARY == warehouse.AutoconfirmIsPrimary {
			centralWarehouses[warehouse.Id] = warehouse
		}

	}

	return &centralWarehouses
}
