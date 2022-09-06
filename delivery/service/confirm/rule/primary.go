package rule

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

// 2. try to confirm from primary warehouse. If 2 warehouses are primary than sort by warehouse.Id asc/desc.
type PrimaryRule struct {
	AbstractRule
}

func (r PrimaryRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	primaryWarehouseIds := *itemConfirm.GetWhRepo().GetCentral()
	if 0 == len(primaryWarehouseIds) {
		itemConfirm.GetLogger().Info("[PrimaryRule] skipped.")
		r.ApplyNextRule(itemConfirm, firstGroup, secondGroup)
		return r
	}

	primaryWarehouseStocks := map[int]int{}
	otherWarehousesStocks := map[int]int{}
	for whId, qty := range *firstGroup {
		if _, ok := primaryWarehouseIds[whId]; ok {
			primaryWarehouseStocks[whId] = qty
		} else {
			otherWarehousesStocks[whId] = qty
		}
	}
	itemConfirm.GetLogger().Info(
		fmt.Sprintf("[PrimaryRule] warehouses %+v.",
			primaryWarehouseStocks,
		))

	if itemConfirm.ComputeConfirmation(&primaryWarehouseStocks, false) {
		return r
	}

	r.ApplyNextRule(itemConfirm, &otherWarehousesStocks, secondGroup)

	return r
}
