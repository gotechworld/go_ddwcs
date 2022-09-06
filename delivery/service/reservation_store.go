package service

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/cms"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/repository"
	atxUtil "gitlab.altex.ro/plug/go_atx_lib/util"
	logger "gitlab.altex.ro/plug/go_logger"
	"strconv"
	"sync"
)

const (
	availableInStoreTrueValue   = 1
	stockInfoFilter             = "info"
	stockWarehouseIsStoreFilter = "warehouse_sync_in_store"
)

type ReservationStoreService struct {
	websiteCode  string
	cmsStoreRepo *repository.StoreRepository
	logger       logger.Logger
}

/**
 * NewReservationStoreService - constructor
 */
func NewReservationStoreService(websiteCode string, logger logger.Logger) *ReservationStoreService {
	repo := repository.NewStoreRepository(
		data_store.CreateCacheCmsStoreInstance(),
		data_store.CreateCmsExternalStore(websiteCode),
		websiteCode,
	)

	return &ReservationStoreService{websiteCode: websiteCode, logger: logger, cmsStoreRepo: repo}
}

/**
 * @param *dto.ReservationStorePayload
 * @return *dto.StoreContainer
 */
func (s *ReservationStoreService) GetAvailableStores(d *dto.ReservationStorePayload) *dto.StoreContainer {
	stocksChannel := make(chan *stocks.StockList)
	storesChannel := make(chan []cms.Store)
	wg := &sync.WaitGroup{}

	wg.Add(2)
	go s.getStockInfo(d, stocksChannel, wg)
	go s.getAvailableCmsStores(storesChannel, wg)
	go func() {
		wg.Wait() // wait until channels values are used
		close(stocksChannel)
		close(storesChannel)
	}()

	return s.getStoresByAvailableWarehouses(s.getAvailableWarehouses(<-stocksChannel, d), <-storesChannel)
}

/**
 * @desc filter available warehouses for reservation to store
 *
 * @param stocksMap *stocks.StockList
 * @param *dto.ReservationStorePayload
 * @return []int
 */
func (s *ReservationStoreService) getAvailableWarehouses(stocksMap *stocks.StockList, d *dto.ReservationStorePayload) []int {
	if stocksMap == nil || len(stocksMap.StocksData.Stocks) == 0 {
		return []int{}
	}

	sellerId := common.GetWebsiteId(d.WebsiteCode)
	warehouses := make([][]int, 0, 1)
	for _, item := range d.Products {
		sellerStock, ok := stocksMap.StocksData.Stocks[item.Sku]
		if !ok { // check for sku in response
			return []int{}
		}

		stockItem, ok := sellerStock[sellerId]
		if !ok || stockItem.Inventory.AvailableInStore != availableInStoreTrueValue {
			// stop execution if there is one product that is not availableInStore - ship to store not possible
			return []int{}
		}

		tmpWarehouses := s.filterWarehousesByQuantity(stockItem.WarehouseInventory, item.Qty)
		if len(tmpWarehouses) == 0 {
			// if there is an empty set, then we won't have any further intersection
			return []int{}
		}

		warehouses = append(warehouses, tmpWarehouses)
	}

	return util.ArrayIntersectMultiple(warehouses...)
}

/**
 * @desc - get store details base on available reservation warehouses
 *
 * @param warehouses []int
 * @param cmsStores []cms.Store
 * @return *dto.StoreContainer
 */
func (s *ReservationStoreService) getStoresByAvailableWarehouses(warehouses []int, cmsStores []cms.Store) *dto.StoreContainer {
	var stores = dto.NewStoreContainer()
	if len(warehouses) == 0 || len(cmsStores) == 0 {
		return stores
	}

	for _, scmStore := range cmsStores {
		if atxUtil.InArray(warehouses, scmStore.ErpID) {
			stores.AddStore(scmStore)
		}
	}

	return stores
}

/**
 * @param warehousesQuantity map[string]interface{}
 * @param quantity int
 * @return []int
 */
func (s *ReservationStoreService) filterWarehousesByQuantity(warehousesQuantity map[string]interface{}, quantity int) []int {
	warehouses := make([]int, 0, 1)
	for warehouse, quantityI := range warehousesQuantity {
		if warehouseQuantity, ok := quantityI.(float64); !ok || int(warehouseQuantity) < quantity {
			continue
		}

		warehouseId, err := strconv.Atoi(warehouse)
		if err != nil {
			s.logger.Error("[ReservationStoreService] [filterWarehousesByQuantity] cannot convert warehouse " + warehouse + "str to int")
			continue
		}

		warehouses = append(warehouses, warehouseId)
	}

	return warehouses
}

/**
 * @desc - get store from external sources
 *
 * @param websiteCode string
 * @return []cms.Store
 */
func (s *ReservationStoreService) getAvailableCmsStores(storesChan chan<- []cms.Store, wg *sync.WaitGroup) {
	//possible filters - commented because we are not 100% sure it should be used
	//"filter1": fmt.Sprintf("network:%s", strconv.Itoa(common.GetWebsiteId(websiteCode))),
	//"filter2": "customer_ship_to_store:1",
	storesChan <- s.cmsStoreRepo.GetStoresBy(map[string]string{})
	wg.Done()
}

/**
 * @param *dto.ReservationStorePayload
 * @return *stocks.StockList
 */
func (s *ReservationStoreService) getStockInfo(d *dto.ReservationStorePayload, channel chan<- *stocks.StockList, wg *sync.WaitGroup) {
	stocksStore := data_store.CreateStocksExternalStore(d.WebsiteCode)
	sellerId := common.GetWebsiteId(d.WebsiteCode)

	params := map[string]string{stockInfoFilter: "3", stockWarehouseIsStoreFilter: "1"}
	for index, item := range d.Products {
		sku := fmt.Sprintf("sku[%d]", index)
		seller := fmt.Sprintf("seller_id[%d]", index)
		params[sku] = item.Sku
		params[seller] = strconv.Itoa(sellerId)
	}

	channel <- stocksStore.GetStocks(params)
	wg.Done()
}
