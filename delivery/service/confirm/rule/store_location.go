package rule

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
	"os"
	"strconv"
	"strings"
)

// 4. This rule enforces confirmation from warehouses that have confirm_regions
// in regions defined in ENV (STORE_LOCATION_RULE_REGIONS) separated by comma (,)
// with a max limit of N qty (defined in GlobalConfig) on each warehouse and applies next rules.
// If finally when all rules have been applied we cannot confirm then
// this rule will remove the max limit of N qty on "targetedWarehouses" warehouses and
// retry confirmation (applies again next rules).
type StoreLocationRule struct {
	AbstractRule
}

func (r StoreLocationRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	limitedQty := r.getQtyLimit()
	regions := r.getRegions()
	if limitedQty <= 0 {
		itemConfirm.GetLogger().Info(fmt.Sprintf(
			"[StoreLocationRule] skipped for sku [%v] due to \"Autoconfirm Store Location QTY Limit\" value: %d",
			regions,
			limitedQty,
		))
		r.ApplyNextRule(itemConfirm, firstGroup, secondGroup)

		return r
	}

	targetedWarehouseIds := *itemConfirm.GetWhRepo().GetStoresInRegions(regions)
	storeLocationLimitedStocks := map[int]int{}
	storeLocationStocks := map[int]int{}
	remotelyWarehousesStocks := map[int]int{}
	for whId, qty := range *firstGroup {
		_, existsTargetedWarehouse := targetedWarehouseIds[whId]
		if existsTargetedWarehouse {
			if limitedQty <= qty {
				storeLocationLimitedStocks[whId] = limitedQty
			} else {
				storeLocationLimitedStocks[whId] = qty
			}

			storeLocationStocks[whId] = qty
		} else {
			remotelyWarehousesStocks[whId] = qty
		}
	}

	// 1. clone and add itemConfirm to a stack for later reuse if cannot confirm
	// 2. try to enforce confirm with QTY limit defined in GlobalConfig and apply next rules
	// 3. if cannot confirm we pop from stack itemConfirm
	// 4. try to enforce confirm without a QTY limit and apply next rules
	itemConfirm.SaveState()
	itemConfirm.GetLogger().Info(fmt.Sprintf(
		"[StoreLocationRule] TRY for region codes %+v with QTY LIMIT [%d] : %v",
		regions,
		limitedQty,
		storeLocationLimitedStocks,
	))

	if itemConfirm.ComputeConfirmation(&storeLocationLimitedStocks, true) {
		itemConfirm.GetLogger().Info(fmt.Sprintf(
			"[StoreLocationRule] for region codes %+v : %v",
			regions,
			storeLocationLimitedStocks,
		))
		return r
	}

	r.ApplyNextRule(itemConfirm, &remotelyWarehousesStocks, secondGroup)

	if !itemConfirm.IsConfirmed() {
		itemConfirm.RestoreState()
	} else {
		return r
	}

	itemConfirm.GetLogger().Info(fmt.Sprintf(
		"[StoreLocationRule] TRY for region codes %+v NO qty limit : %+v",
		regions,
		storeLocationStocks,
	))

	if itemConfirm.ComputeConfirmation(&storeLocationStocks, true) {
		itemConfirm.GetLogger().Info(fmt.Sprintf("[StoreLocationRule] warehouses : %+v", storeLocationStocks))

		return r
	}

	r.ApplyNextRule(itemConfirm, &remotelyWarehousesStocks, secondGroup)

	return r
}

func (r StoreLocationRule) getQtyLimit() int {
	limit, _ := strconv.Atoi(os.Getenv("DELIVERY_STORE_LOCATION_QTY_LIMIT"))
	return limit
}

func (r StoreLocationRule) getRegions() []string {
	return strings.Split(os.Getenv("STORE_LOCATION_RULE_REGIONS"), ",")
}
