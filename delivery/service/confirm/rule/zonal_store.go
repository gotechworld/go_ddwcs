package rule

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

// 3. enforce confirmation from Zonal Stores (same confirm_region with customer shipping region).
type ZonalStoreRule struct {
	AbstractRule
}

func (r ZonalStoreRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	storesInRegionIds := *itemConfirm.GetWhRepo().GetStoresInRegions([]string{itemConfirm.GetRegionCode()})
	if 0 == len(storesInRegionIds) {
		itemConfirm.GetLogger().Info("[ZonalStoreRule] skipped.")
		r.ApplyNextRule(itemConfirm, firstGroup, secondGroup)
		return r
	}

	storesInRegionStocks := map[int]int{}
	otherStocks := map[int]int{}
	for whId, qty := range *firstGroup {
		if _, ok := storesInRegionIds[whId]; ok {
			storesInRegionStocks[whId] = qty
		} else {
			otherStocks[whId] = qty
		}
	}

	if itemConfirm.ComputeConfirmation(&storesInRegionStocks, true) {
		itemConfirm.GetLogger().Info(
			fmt.Sprintf("[ZonalStoreRule] warehouses %+v.",
				storesInRegionStocks,
			))
		return r
	}

	r.ApplyNextRule(itemConfirm, &otherStocks, secondGroup)

	return r
}
