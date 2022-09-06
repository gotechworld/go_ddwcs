package rule

import (
	"fmt"
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
	"sort"
)

// 5. Confirm from warehouses with low priority first. enforce confirmation for warehouses with priority <= 10
type PriorityRule struct {
	AbstractRule
}

func (r PriorityRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	if !r.CanApply(itemConfirm, firstGroup) {
		return r
	}

	noThresholdWarehouses := *itemConfirm.GetWhRepo().GetAvailable()
	allEnabledWarehouses := *itemConfirm.GetWhRepo().GetAll()
	priorityWarehousesGroups := map[int]map[int]int{}
	for whId, qty := range *firstGroup {
		whPriority := allEnabledWarehouses[whId].Priority
		// remove from G1 Warehouses that reached defined threshold
		if r.enforceConfirmationForPriorityLevel(whPriority) {
			if _, ok := noThresholdWarehouses[whId]; !ok {
				itemConfirm.GetLogger().Info(
					fmt.Sprintf(
						"[PriorityRule] G1 ( priority <= 10 ) whId [%d] removed because today confirmed products >= wh threshold.",
						whId,
					))
				continue
			}
		}

		priorityWarehousesGroup, ok := priorityWarehousesGroups[whPriority]
		if ok {
			priorityWarehousesGroup[whId] = qty
		} else {
			priorityWarehousesGroup = map[int]int{whId: qty}
		}
		priorityWarehousesGroups[whPriority] = priorityWarehousesGroup
	}

	if 0 == len(priorityWarehousesGroups) {
		r.ApplyNextRule(itemConfirm, firstGroup, secondGroup)

		return r
	}

	priorities := make([]int, 0, len(priorityWarehousesGroups))
	for priority := range priorityWarehousesGroups {
		priorities = append(priorities, priority)
	}
	sort.Ints(priorities)

	itemConfirm.GetLogger().Info(fmt.Sprintf("[PriorityRule] warehouses : %+v", priorityWarehousesGroups))

	for _, priorityLevel := range priorities {
		priorityWarehousesGroup := priorityWarehousesGroups[priorityLevel]
		if 1 == len(priorityWarehousesGroup) {
			if itemConfirm.ComputeConfirmation(&priorityWarehousesGroup, r.enforceConfirmationForPriorityLevel(priorityLevel)) {
				return r
			}

			continue
		}

		r.ApplyNextRule(itemConfirm, &priorityWarehousesGroup, secondGroup)
	}

	return r
}

func (r PriorityRule) enforceConfirmationForPriorityLevel(priority int) bool {
	return priority <= 10
}
