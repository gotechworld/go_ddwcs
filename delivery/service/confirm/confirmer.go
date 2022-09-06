package confirm

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/common/dto/stocks"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/dto"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/repository"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/item_confirm"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/rule"
	logger "gitlab.altex.ro/plug/go_logger"
	"reflect"
	"strconv"
)

type Confirmer struct {
	whRepo        *repository.WarehouseRepository
	logger        logger.Logger
	stocksMap     *map[string]stocks.SellerStock
	confirmations *map[string]map[int]int
}

// Confirmer - constructor
func NewConfirmer(stocksMap *map[string]stocks.SellerStock, whRepo *repository.WarehouseRepository, logger logger.Logger) *Confirmer {
	return &Confirmer{
		whRepo:        whRepo,
		logger:        logger,
		stocksMap:     stocksMap,
		confirmations: &map[string]map[int]int{},
	}
}

// Confirm given products. To get confirmed product data use GetConfirmations() method.
//
// @param products *dto.EstimateProducts - ["productSku": map[string]dto.ProductInfo, ...]
// @param deliveryRegion string
//
func (c *Confirmer) ConfirmItems(products *dto.EstimateProducts, deliveryRegion string) {
	ruleChain := c.getRuleChain()
	// confirm each product
	for sku, _ := range *products {
		itemStocks := c.getItemStocks(sku)
		itemConfirm := item_confirm.NewItemConfirm(c.whRepo, c.logger, (*products)[sku].Qty, deliveryRegion)
		ruleChain.Apply(itemConfirm, itemStocks, &map[int]int{})
		if itemConfirm.IsConfirmed() {
			if 0 == len(itemConfirm.Confirmations) {
				(*c.confirmations)[sku] = itemConfirm.ForcedConfirmations
				c.logger.Info(
					fmt.Sprintf("[ForcedConfirmations] for [%s] : %+v.",
						sku,
						itemConfirm.ForcedConfirmations,
					))
			} else {
				(*c.confirmations)[sku] = itemConfirm.Confirmations
				c.logger.Info(
					fmt.Sprintf("[Confirmations] for [%s] : %+v.",
						sku,
						itemConfirm.Confirmations,
					))
			}
		} else {
			c.logger.Info(
				fmt.Sprintf("[NO Confirmations] for [%s] : []",
					sku,
				))
		}
	}
}

// Build a chain of Rules that will be used to confirm products(items).
//
// @return rule.Rule - first rule in the chain of Rules
//
func (c *Confirmer) getRuleChain() rule.Rule {
	zonalRule := rule.ZonalRule{}
	primaryRule := rule.PrimaryRule{}
	zonalStoreRule := rule.ZonalStoreRule{}
	storeLocationRule := rule.StoreLocationRule{}
	priorityRule := rule.PriorityRule{}
	maxStockRule := rule.MaxStockRule{}
	proximityRule := rule.ProximityRule{}
	consolidationRule := rule.ConsolidationRule{}
	thresholdRule := rule.ThresholdRule{}

	consolidationRule.NextRule = thresholdRule
	proximityRule.NextRule = consolidationRule
	maxStockRule.NextRule = proximityRule
	priorityRule.NextRule = maxStockRule
	storeLocationRule.NextRule = priorityRule
	zonalStoreRule.NextRule = storeLocationRule
	primaryRule.NextRule = zonalStoreRule
	zonalRule.NextRule = primaryRule

	return zonalRule
}

// Get stocks for given product sku.
// Supplier stocks are aggregated.
//
// @param sku string - a product sku
// @return *map[int]int - map[whId: qty, ...]
//
func (c *Confirmer) getItemStocks(sku string) *map[int]int {
	itemStocks := map[int]int{}
	for _, stockData := range (*c.stocksMap)[sku] {
		for whId, stockData := range stockData.WarehouseInventory {
			whIdInt, _ := strconv.Atoi(whId)
			rt := reflect.TypeOf(stockData)
			qty := 0
			if rt.Kind() == reflect.Array {
				data := stockData.([]string)
				for _, supplierQty := range data {
					supplierQtyInt, _ := strconv.Atoi(supplierQty)
					qty = qty + supplierQtyInt
				}
			} else {
				qty, _ = strconv.Atoi(fmt.Sprintf("%v", stockData))
			}

			itemStocks[whIdInt] = qty
		}
	}

	return &itemStocks
}

// Return Confirmations computed by ConfirmItems() method.
//
// @return *map[string]map[int]int - map[whId: Qty], ...]
//
func (c *Confirmer) GetConfirmations() *map[string]map[int]int {
	return c.confirmations
}
