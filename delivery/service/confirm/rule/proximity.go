package rule

import (
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

// 7. Confirm from warehouses that are in same region with suborder address region (delivery region).
type ProximityRule struct {
	AbstractRule
}

func (r ProximityRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	r.ApplyNextRule(itemConfirm, firstGroup, secondGroup)

	return r
}
