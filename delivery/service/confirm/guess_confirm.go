package confirm

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/repository"
	logger "gitlab.altex.ro/plug/go_logger"
	"os"
	"strconv"
)

type GuessConfirm struct {
	whRepo *repository.WarehouseRepository
	logger logger.Logger
}

// GuessConfirm - constructor
func NewGuessConfirm(whRepo *repository.WarehouseRepository, logger logger.Logger) *GuessConfirm {
	return &GuessConfirm{whRepo: whRepo, logger: logger}
}

// Try to guess confirm warehouse group and return group "delivery time" in hours (zonal/central/supplier/default).
// @TODO : implement confirm logic from OMS in item_confirm.go
func (gc GuessConfirm) GetEstimate(stocks *map[string]stocks.SellerStock, regionCode string) int {
	// 1. check if we have at leas one product ONLY in Supplier
	if gc.haveProductOnlyInSupplierStock(stocks) {
		supplierDeliveryTime, _ := strconv.Atoi(os.Getenv("DELIVERY_TIME_SUPPLIER"))
		gc.logger.Printf(
			"[Estimator] [can deliver from SUPPLIER] returned supplier delivery time [%+v] days.",
			supplierDeliveryTime/24,
		)
		return supplierDeliveryTime
	}

	// 2. check if we can confirm from Zonal Warehouse
	// (enforce = confirm from zonal as much as we can - but on guess we don't have qty's)
	if gc.canConfirmFromZonal(stocks, regionCode) {
		zonalDeliveryTime, _ := strconv.Atoi(os.Getenv("DELIVERY_TIME_ZONAL"))
		gc.logger.Printf(
			"[Estimator] [can deliver from ZONAL] returned zonal delivery time [%+v] days.",
			zonalDeliveryTime/24,
		)
		return zonalDeliveryTime
	}

	// 3. check if we can confirm from Primary Warehouse
	// (all qty or nothing - but on guess we don't have qty's)
	if gc.canConfirmFromCentral(stocks) {
		centralDeliveryTime, _ := strconv.Atoi(os.Getenv("DELIVERY_TIME_CENTRAL"))
		gc.logger.Printf(
			"[Estimator] [can deliver from CENTRAL] returned central delivery time [%+v] days.",
			centralDeliveryTime/24,
		)
		return centralDeliveryTime
	}

	// 4. guess confirm based on priority => DELIVERY_TIME_DEFAULT
	defaultDeliveryTime, _ := strconv.Atoi(os.Getenv("DELIVERY_TIME_DEFAULT"))
	gc.logger.Printf(
		"[Estimator] [DEFAULT] returned default delivery time [%+v] days.",
		defaultDeliveryTime/24,
	)

	return defaultDeliveryTime
}

// @TODO move in zonal rule
func (gc GuessConfirm) canConfirmFromZonal(stocks *map[string]stocks.SellerStock, regionCode string) bool {
	if regionCode == "" {
		return false
	}

	zonalWarehouses := *gc.whRepo.GetZonal(regionCode)
	if 0 == len(zonalWarehouses) {
		return false
	}

	canConfirmFromZonal := map[string]bool{}
	for productSku, sellerStock := range *stocks {
		canConfirmFromZonal[productSku] = false
		for _, stock := range sellerStock {
			for warehouseId := range stock.WarehouseInventory {
				whId, _ := strconv.Atoi(warehouseId)
				_, exist := zonalWarehouses[whId]
				if exist {
					canConfirmFromZonal[productSku] = true
					break
				}
			}
		}
	}

	for _, canConfirm := range canConfirmFromZonal {
		if !canConfirm {
			return false
		}
	}

	return true
}

// @TODO move in primary rule
func (gc GuessConfirm) canConfirmFromCentral(stocks *map[string]stocks.SellerStock) bool {
	centralWarehouses := *gc.whRepo.GetCentral()
	if 0 == len(centralWarehouses) {
		return false
	}

	canConfirmFromCentral := map[string]bool{}
	for productSku, sellerStock := range *stocks {
		canConfirmFromCentral[productSku] = false
		for _, stock := range sellerStock {
			for warehouseId := range stock.WarehouseInventory {
				whId, _ := strconv.Atoi(warehouseId)
				_, exist := centralWarehouses[whId]
				if exist {
					canConfirmFromCentral[productSku] = true
					break
				}
			}
		}
	}

	for _, canConfirm := range canConfirmFromCentral {
		if !canConfirm {
			return false
		}
	}

	return true
}

// check if we have at least one product that is in ONLY in Supplier Stock
// than all products will have Supplier's delivery time
func (gc GuessConfirm) haveProductOnlyInSupplierStock(stocks *map[string]stocks.SellerStock) bool {
	suppliersWarehouseId := os.Getenv("WAREHOUSE_ID_FOR_SUPPLIERS")

	for _, sellerStock := range *stocks {
		for _, stock := range sellerStock {
			warehousesNr := len(stock.WarehouseInventory)
			// if more than 1 means we can confirm from other warehouse than supplier
			if 1 < warehousesNr {
				continue
			}

			for warehouseId := range stock.WarehouseInventory {
				if suppliersWarehouseId == warehouseId {
					return true
				}
			}
		}
	}

	return false
}
