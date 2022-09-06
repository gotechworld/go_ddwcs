package stocks

import (
	"encoding/json"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/net/rest"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
	"strconv"
)

const STOCKS_API_WAREHOUSE_ENDPOINT = "/api/v1.0/warehouses"

type ExternalWarehousesStore struct {
	stocksClient *rest.StocksClient
	logger.LoggerAware
}

// ExternalWarehousesStore - constructor
func NewExternalWarehousesStore(stocksClient *rest.StocksClient) *ExternalWarehousesStore {
	return &ExternalWarehousesStore{stocksClient: stocksClient}
}

// Get Warehouses List
// @param params map[string]string
// @return []stocks.Warehouse
// @desc - get all warehouses from Stock API
// the params map can contain the following keys: // @TODO ... I would transform this in a struct
func (ews *ExternalWarehousesStore) GetWarehouses(params map[string]string) []stocks.Warehouse {
	var warehouses []stocks.Warehouse
	if ews.stocksClient == nil {
		ews.GetLogger().Printf("[ERROR] [StocksAPI]: 'stocksClient' is undefined")
		return warehouses
	}

	currentPage := 1
	for {
		params["page_no"] = strconv.Itoa(currentPage)
		responseBytes, err := ews.stocksClient.Get(STOCKS_API_WAREHOUSE_ENDPOINT, params)
		if err != nil {
			ews.GetLogger().Printf(
				"[ERROR] [StocksAPI] [%s]: %+v",
				ews.stocksClient.GetUrl(),
				err,
			)
			return warehouses
		}

		var warehousesResponse stocks.WarehousesList
		err = json.Unmarshal(responseBytes, &warehousesResponse)
		if err != nil {
			ews.GetLogger().Printf(
				"[ERROR] [StocksAPI] [%s] on Unmarshal: %s. Response unmarshalled: %s",
				ews.stocksClient.GetUrl(),
				err,
				responseBytes,
			)
			return warehouses
		}

		ews.GetLogger().Printf(
			"[StocksAPI] [%s]",
			ews.stocksClient.GetUrl(),
		)
		warehouses = append(warehouses, warehousesResponse.WarehousesData.Warehouses...)
		if warehousesResponse.WarehousesData.TotalPages == warehousesResponse.WarehousesData.CurrentPage {
			break
		}

		currentPage++
	}

	return warehouses
}
