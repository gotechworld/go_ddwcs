package rule

import (
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

// 8. Confirm from warehouse that have "stock qty" = (needed - confirmed).
type ConsolidationRule struct {
	AbstractRule
}

func (r ConsolidationRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	r.ApplyNextRule(itemConfirm, firstGroup, secondGroup)

	return r
}
