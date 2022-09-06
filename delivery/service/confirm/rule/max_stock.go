package rule

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
	"sort"
)

// 6. Confirm from Warehouses with bigger stock qty's.
type MaxStockRule struct {
	AbstractRule
}

func (r MaxStockRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	if !r.CanApply(itemConfirm, firstGroup) {
		return r
	}

	// group and sort by stock qty
	sortedByQtyWarehousesGroups := map[int]map[int]int{}
	for whId, qty := range *firstGroup {
		sortedByQtyWarehouses, ok := sortedByQtyWarehousesGroups[qty]
		if ok {
			sortedByQtyWarehouses[whId] = qty
		} else {
			sortedByQtyWarehouses = map[int]int{whId: qty}
		}
		sortedByQtyWarehousesGroups[qty] = sortedByQtyWarehouses
	}

	quantities := make([]int, 0, len(sortedByQtyWarehousesGroups))
	for qty := range sortedByQtyWarehousesGroups {
		quantities = append(quantities, qty)
	}
	sort.Ints(quantities)

	// apply next rule for groups with same qty
	for _, qty := range quantities {
		sortedByQtyWarehousesGroup := sortedByQtyWarehousesGroups[qty]
		if 1 == len(sortedByQtyWarehousesGroup) {
			if itemConfirm.ComputeConfirmation(&sortedByQtyWarehousesGroup, false) {
				itemConfirm.GetLogger().Info(fmt.Sprintf("[MaxStockRule] warehouses : %+v", sortedByQtyWarehousesGroup))
				return r
			}

			continue
		}

		r.ApplyNextRule(itemConfirm, &sortedByQtyWarehousesGroup, secondGroup)
	}

	return r
}
