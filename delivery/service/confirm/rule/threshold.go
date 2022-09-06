package rule

import (
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

// 9. Confirm from warehouses that have not reach specified Threshold for current day.
type ThresholdRule struct {
	AbstractRule
}

func (r ThresholdRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	if !r.CanApply(itemConfirm, firstGroup) {
		return r
	}

	return r
}
