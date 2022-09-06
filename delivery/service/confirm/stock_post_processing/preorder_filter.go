package stock_post_processing

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	logger "gitlab.altex.ro/plug/go_logger"
)

type PreorderFilter struct {
	logger logger.Logger
}

// PreorderFilter - constructor
func NewPreorderFilter(logger logger.Logger) *PreorderFilter {
	preorderFilter := &PreorderFilter{
		logger: logger,
	}

	return preorderFilter
}

// Remove preorder products from stocksMap
func (pf PreorderFilter) Apply(stocksMap *map[string]stocks.SellerStock) {
	for productSku, sellerStock := range *stocksMap {
		for _, stockData := range sellerStock {
			if stocks.STATUS_PREORDER == stockData.Inventory.Status {
				delete(*stocksMap, productSku)
				pf.logger.Info(
					fmt.Sprintf("[PreorderFilter] [remove] preorder product [%+v].",
						productSku,
					))
			}
		}
	}
}
