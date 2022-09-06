package stocks

import (
	"encoding/json"
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/net/rest"
	"gitlab.altex.ro/plug/go_atx_lib/logger"
	"strconv"
)

const STOCKS_API_STOCK_ENDPOINT = "/api/v1.0/stocks"

type ExternalStocksStore struct {
	stocksClient *rest.StocksClient
	logger.LoggerAware
}

// ExternalStocksStore - constructor
func NewExternalStocksStore(stocksClient *rest.StocksClient) *ExternalStocksStore {
	return &ExternalStocksStore{stocksClient: stocksClient}
}

// Get stock information for specified skus/seller_id pairs.
//
// @param params map[string]string
// @return stocks.StockList
// @desc - get the stock values from Stock API
// the params map will contain the following keys: "skus", "seller_id" for each skus
func (ess *ExternalStocksStore) GetStocks(params map[string]string) *stocks.StockList {
	emptyList := &stocks.StockList{}
	if len(params) == 0 {
		ess.GetLogger().Print("[StocksAPI]: empty sku list")
		return emptyList
	}

	if nil == ess.stocksClient {
		ess.GetLogger().Error("[StocksAPI]: 'stocksClient' is undefined")
		return emptyList
	}

	responseBytes, err := ess.stocksClient.Get(STOCKS_API_STOCK_ENDPOINT, params)

	if err != nil {
		ess.GetLogger().Printf("[ERROR] [StocksAPI] [%s] : %+v", ess.stocksClient.GetUrl(), err)
		return emptyList
	}

	var stockResponse *stocks.StockList
	err = json.Unmarshal(responseBytes, &stockResponse)
	if err != nil {
		ess.GetLogger().Printf(
			"[ERROR] [StocksAPI] [%s] on Unmarshal: %s. Response unmarshalled: %s",
			ess.stocksClient.GetUrl(),
			err,
			responseBytes,
		)
		return emptyList
	}

	ess.GetLogger().Printf("[StocksAPI] [%s] ", ess.stocksClient.GetUrl())

	return stockResponse
}

// Prepare GetStocks params if necessary adding seller_id[%key]= sellerId params and formatting sku[%key] = sku
func (ess *ExternalStocksStore) PrepareParams(sku []string, sellerId int) map[string]string {
	params := map[string]string{}
	for index, productSku := range sku {
		paramKey := fmt.Sprintf("sku[%v]", index)
		params[paramKey] = productSku
		paramKey = fmt.Sprintf("seller_id[%v]", index)
		params[paramKey] = strconv.Itoa(sellerId)
	}

	return params
}
