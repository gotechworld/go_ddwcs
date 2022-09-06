package rule

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

// 1. enforce confirmation from Zonal Warehouses (same confirm_region with customer shipping region).
type ZonalRule struct {
	AbstractRule
}

func (r ZonalRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	zonalWarehouseIds := *itemConfirm.GetWhRepo().GetZonal(itemConfirm.GetRegionCode())
	if 0 == len(zonalWarehouseIds) {
		itemConfirm.GetLogger().Info("[ZonalRule] skipped.")
		r.ApplyNextRule(itemConfirm, firstGroup, secondGroup)
		return r
	}

	zonalWarehouseStocks := map[int]int{}
	remotelyWarehousesStocks := map[int]int{}
	for whId, qty := range *firstGroup {
		if _, ok := zonalWarehouseIds[whId]; ok {
			zonalWarehouseStocks[whId] = qty
		} else {
			remotelyWarehousesStocks[whId] = qty
		}
	}

	if itemConfirm.ComputeConfirmation(&zonalWarehouseStocks, true) {
		itemConfirm.GetLogger().Info(
			fmt.Sprintf("[ZonalRule] warehouses %+v.",
				zonalWarehouseStocks,
			))
		return r
	}

	r.ApplyNextRule(itemConfirm, &remotelyWarehousesStocks, secondGroup)

	return r
}
