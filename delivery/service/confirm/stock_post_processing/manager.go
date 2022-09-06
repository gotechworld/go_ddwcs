package stock_post_processing

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	logger "gitlab.altex.ro/plug/go_logger"
)

type StockPostProcessing interface {
	Apply(stocksMap *map[string]stocks.SellerStock)
}

type Manager struct {
	stockPostProcessors map[string]StockPostProcessing
}

func NewManager(logger logger.Logger) *Manager {
	return &Manager{stockPostProcessors: map[string]StockPostProcessing{
		"preorder_filter": NewPreorderFilter(logger),
	}}
}

func (m Manager) Apply(stocksMap *map[string]stocks.SellerStock) {
	for _, stockProcessor := range m.stockPostProcessors {
		stockProcessor.Apply(stocksMap)
	}
}
