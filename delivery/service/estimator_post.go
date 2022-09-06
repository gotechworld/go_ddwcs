package service

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	ddwcsUtil "gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/repository"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/stock_post_processing"
	"gitlab.altex.ro/plug/go_atx_lib/util"
	logger "gitlab.altex.ro/plug/go_logger"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type DeliveryPostEstimator struct {
	requestData *dto.EstimateRequestPost
	whRepo      *repository.WarehouseRepository
	logger      logger.Logger
}

// NewDeliveryPostEstimator - DeliveryPostEstimator constructor
func NewDeliveryPostEstimator(logger logger.Logger, requestData *dto.EstimateRequestPost) *DeliveryPostEstimator {
	whRepo := repository.NewWarehouseRepository(requestData.WebsiteCode, logger)
	estimator := &DeliveryPostEstimator{
		requestData: requestData,
		logger:      logger,
		whRepo:      whRepo,
	}

	return estimator
}

// EstimateDelivery - Computes nr of estimated days based on confirmation algorithm and mapped interval for each warehouse.
//
// @return (response *dto.EstimateResponsePost, statusCode int)
//
func (de DeliveryPostEstimator) EstimateDelivery() (response *dto.EstimateResponsePost, statusCode int) {
	stocksMap := de.getStocksFromStocksAPI()
	// 1. confirm each product
	confirmations := de.ComputeConfirmations(stocksMap)
	// 2. get carrier and estimate and build the response
	response = de.BuildEstimateResponse(confirmations)
	// 3. enrich response with "generic" key
	de.AddGenericKey(response, stocksMap)

	return response, http.StatusOK
}

func (de DeliveryPostEstimator) getStocksFromStocksAPI() map[string]stocks.SellerStock {
	skus := make([]string, 0, len(de.requestData.Products))
	for sku, productInfo := range de.requestData.Products {
		if 1 == productInfo.IsPreOrder {
			continue
		}

		skus = append(skus, sku)
	}
	stocksStore := data_store.CreateStocksExternalStore(de.requestData.WebsiteCode)
	params := stocksStore.PrepareParams(skus, common.GetWebsiteId(de.requestData.WebsiteCode))

	return stocksStore.GetStocks(params).StocksData.Stocks
}

// ComputeConfirmations - Compute confirmation for given product skus.
//
// @param stocksMap map[string]stocks.SellerStock
// @return*map[string]map[int]int - map["product_sku": map [whId: qty, ...], ...]
//
func (de DeliveryPostEstimator) ComputeConfirmations(stocksMap map[string]stocks.SellerStock) *map[string]map[int]int {
	stockPostProcessingManager := stock_post_processing.NewManager(de.logger)
	stockPostProcessingManager.Apply(&stocksMap)
	confirmManger := confirm.NewConfirmer(&stocksMap, de.whRepo, de.logger)
	confirmManger.ConfirmItems(&de.requestData.Products, de.requestData.DeliveryRegion)

	return confirmManger.GetConfirmations()
}

// BuildEstimateResponse - Build estimate response based on given confirmation map
//
// @param confirmations *map[string]map[int]int - map["product_sku": map [whId: qty, ...], ...]
// @return *dto.EstimateResponsePost
//
func (de DeliveryPostEstimator) BuildEstimateResponse(confirmations *map[string]map[int]int) *dto.EstimateResponsePost {
	allWarehouses := *de.whRepo.GetAll()
	response := &dto.EstimateResponsePost{}
	estimatedProductInfo := map[string]dto.EstimatedProductInfo{}
	for sku, confirmation := range *confirmations {
		attributeSetId := de.requestData.Products[sku].AttributeSetId
		for whId, _ := range confirmation {
			carrier := allWarehouses[whId].DefaultCourier
			for _, courier := range allWarehouses[whId].Couriers {
				if util.InArray(courier.AttributeSets, attributeSetId) {
					carrier = courier.Name
				}
			}

			estimatedProductInfo[sku] = dto.EstimatedProductInfo{
				BundleId:              de.requestData.Products[sku].BundleId,
				EstimatedDeliveryDays: allWarehouses[whId].DeliveryEstimate,
				DeliverFromRegion:     allWarehouses[whId].RegionName,
				Carrier:               carrier,
			}
		}
	}

	// populate response with skus that cannot be confirmed
	for sku, productInfo := range de.requestData.Products {
		if _, ok := estimatedProductInfo[sku]; ok {
			continue
		}

		estimatedProductInfo[sku] = dto.EstimatedProductInfo{
			BundleId:              productInfo.BundleId,
			EstimatedDeliveryDays: os.Getenv("DELIVERY_TIME_NO_STOCK_INTERVAL"),
			DeliverFromRegion:     os.Getenv("DELIVERY_TIME_NO_STOCK_REGION"),
			Carrier:               os.Getenv("DELIVERY_TIME_NO_STOCK_CARRIER"),
		}
	}

	response.EstimatedDelivery = estimatedProductInfo

	return response
}

// AddGenericKey - Build and add "generic" response key based on given estimateResponsePost.
// 1. compute and add EstimatedDelivery interval
// 2. compute and add EstimatedDeliveryDates
// 3. compute and add AdditionalDeliveryDates
// 4. compute and add ShipToStore
// @param stocksMap map[string]stocks.SellerStock
// @param response *dto.EstimateResponsePost
//
func (de DeliveryPostEstimator) AddGenericKey(response *dto.EstimateResponsePost, stocksMap map[string]stocks.SellerStock) {
	// 1. compute and add "Estimated Delivery Interval"
	minDeliveryDays := math.MaxInt16
	maxDeliveryDays := 0
	for _, estimatedDeliveryInfo := range response.EstimatedDelivery {
		estimatedDays := strings.Split(estimatedDeliveryInfo.EstimatedDeliveryDays, "-")
		for i := range estimatedDays {
			estimatedDays[i] = strings.TrimSpace(estimatedDays[i])
		}
		startDay, _ := strconv.Atoi(estimatedDays[0])
		endDay, _ := strconv.Atoi(estimatedDays[1])

		if startDay < minDeliveryDays {
			minDeliveryDays = startDay
		}

		if endDay > maxDeliveryDays {
			maxDeliveryDays = endDay
		}
	}
	response.Generic.EstimatedDelivery = strconv.Itoa(minDeliveryDays) + "-" + strconv.Itoa(maxDeliveryDays)
	// 2. compute and add "Estimated Delivery Dates"
	nationalHolidays := repository.NewGlobalConfigRepository().GetNationalHolidays(de.requestData.WebsiteCode)
	response.Generic.EstimatedDeliveryDates = de.generateDatesInterval(minDeliveryDays, maxDeliveryDays, nationalHolidays)
	// 3. compute and add "Additional Delivery Dates"
	additionalStartDay := maxDeliveryDays + 1
	additionalDeliveryDays, _ := strconv.Atoi(os.Getenv("ADDITIONAL_DELIVERY_DAYS"))
	additionalEndDay := additionalStartDay + additionalDeliveryDays
	response.Generic.AdditionalDeliveryDates = de.generateDatesInterval(additionalStartDay, additionalEndDay, nationalHolidays)
	// 4. compute and add "ship_to_store" key
	response.Generic.ShipToStore = 0
	if de.isShipToStoreAllowed(stocksMap) {
		response.Generic.ShipToStore = 1
	}
}

// Generates dates from current date in [currentDate + startDay; currentDate + endDay] interval
//
// @param int startDay
// @param int endDay
// @param *map[string]string nationalHolidays {"YYYY-MM-DD":"YYYY-MM-DD", ...}
// @return []string - array of date strings in "YYYY-MM-DD" format
//
func (de DeliveryPostEstimator) generateDatesInterval(startDay int, endDay int, nationalHolidays *map[string]string) []string {
	currentTime := ddwcsUtil.GetLocalTime(time.Now())
	var estimatedDeliveryDates []string
	for i := startDay; i <= endDay; i++ {
		dateToAdd := currentTime.AddDate(0, 0, i)
		dateToAddWeekday := dateToAdd.Weekday()
		if time.Sunday == dateToAddWeekday {
			continue
		}

		_, isNationalHoliday := (*nationalHolidays)[dateToAdd.Format("2006-01-02")]
		if isNationalHoliday {
			continue
		}

		estimatedDeliveryDates = append(estimatedDeliveryDates, dateToAdd.Format("2006-01-02"))
	}

	return estimatedDeliveryDates
}

// Computes if current products can have ship_to_store
// true - if all products have enough stocks in Central WHs (Primary & enabled), false - otherwise
// @param stocksMap map[string]stocks.SellerStock
// @return bool
func (de DeliveryPostEstimator) isShipToStoreAllowed(stocksMap map[string]stocks.SellerStock) bool {
	centralWarehouses := *de.whRepo.GetCentral()//
	for sku, productInfo := range de.requestData.Products {
		if 1 == productInfo.IsPreOrder {
			continue
		}

		currentSkuHasStockInCentralWh := false
		for _, stockData := range stocksMap[sku] {
			for whId, stockData := range stockData.WarehouseInventory {
				whIdInt, _ := strconv.Atoi(whId)
				whQty := 0
				_, isCentralWh := centralWarehouses[whIdInt]
				if !isCentralWh {
					continue
				}

				rt := reflect.TypeOf(stockData)
				if rt.Kind() == reflect.Array {
					data := stockData.([]string)
					for _, supplierQty := range data {
						supplierQtyInt, _ := strconv.Atoi(supplierQty)
						whQty = whQty + supplierQtyInt
					}
				} else {
					whQty, _ = strconv.Atoi(fmt.Sprintf("%v", stockData))
				}
				if whQty >= productInfo.Qty {
					currentSkuHasStockInCentralWh = true
				}
			}
		}

		// if one product don't have enough stock a Central Wh ship_to_store is NOT allowed
		if !currentSkuHasStockInCentralWh {
			return false
		}
	}

	// all products have enough stock in Central Whs so ship_to_store is allowed
	return true
}
