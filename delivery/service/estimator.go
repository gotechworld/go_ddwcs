package service

import (
	"gitlab.altex.ro/ams/go_ddwcs/common"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/data_store"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/repository"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/stock_post_processing"
	logger "gitlab.altex.ro/plug/go_logger"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DeliveryEstimator struct {
	websiteCode    string
	confirmGuesser *confirm.GuessConfirm
	logger         logger.Logger
}

// DeliveryEstimator - constructor
func NewDeliveryEstimator(logger logger.Logger, websiteCode string) *DeliveryEstimator {
	whRepo := repository.NewWarehouseRepository(websiteCode, logger)
	estimator := &DeliveryEstimator{
		websiteCode:    websiteCode,
		logger:         logger,
		confirmGuesser: confirm.NewGuessConfirm(whRepo, logger),
	}

	return estimator
}

// Computes nr of estimated day based on guess confirmation and mapped hours for each warehouse category
// returns "delivery estimate days" and "http statusCode"
func (de DeliveryEstimator) EstimateDelivery(sku []string, regionCode string) (estimatedDays int, statusCode int) {
	stocksStore := data_store.CreateStocksExternalStore(de.websiteCode)
	params := stocksStore.PrepareParams(sku, common.GetWebsiteId(de.websiteCode))
	stocksMap := stocksStore.GetStocks(params).StocksData.Stocks
	stockPostProcessingManager := stock_post_processing.NewManager(de.logger)
	stockPostProcessingManager.Apply(&stocksMap)

	if 0 == len(stocksMap) || !de.allProductsHaveStocks(&stocksMap, sku) {
		defaultDeliveryTime, _ := strconv.Atoi(os.Getenv("DELIVERY_TIME_NO_STOCK"))
		defaultDeliveryTime /= 24
		de.logger.Printf(
			"[Estimator] [NO Available Stocks] returned DELIVERY_TIME_NO_STOCK [%+v] days.",
			defaultDeliveryTime,
		)

		return defaultDeliveryTime, http.StatusNotFound
	}

	estimatedHours := de.confirmGuesser.GetEstimate(&stocksMap, regionCode)

	estimatedHours = de.applyCorrections(estimatedHours)
	de.logger.Printf(
		"[Estimator] [Correction] after apply correction time [%+v] days.",
		estimatedHours/24,
	)

	return estimatedHours / 24, http.StatusOK
}

// Check if all product skus are in stocks API.
func (de DeliveryEstimator) allProductsHaveStocks(stocks *map[string]stocks.SellerStock, skus []string) bool {
	for _, productSku := range skus {
		_, contains := (*stocks)[productSku]
		if !contains {
			return false
		}
	}

	return true
}

// Apply hour and weekdays corrections
func (de DeliveryEstimator) applyCorrections(hours int) int {
	return de.applyCorrectionByWeekDay(de.applyCorrectionByHour(hours))
}

// Apply Hour Correction.
// All orders placed after 12 will take at least one extra-day
// Don't apply this rule for Friday, Saturday and Sunday
func (de DeliveryEstimator) applyCorrectionByHour(hours int) int {
	currentTime := util.GetLocalTime(time.Now())
	currentWeekday := currentTime.Weekday()
	if time.Saturday == currentWeekday || time.Sunday == currentWeekday || time.Friday == currentWeekday {
		return hours
	}

	if 12 < currentTime.Hour() {
		return hours + 24
	}

	return hours
}

// Apply correction based on current weekday.
// Through weekends, central/zonal warehouses are working at a lower rate (50%)
// and store warehouses will not pack products at all
func (de DeliveryEstimator) applyCorrectionByWeekDay(hours int) int {
	currentWeekday := util.GetLocalTime(time.Now()).Weekday()
	switch {
	case time.Saturday == currentWeekday:
		return hours + 24
	case time.Sunday == currentWeekday:
		return hours + 24
	case time.Friday == currentWeekday:
		return hours + 48
	}

	return hours
}
