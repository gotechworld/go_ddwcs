package rule

import (
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

type AbstractRule struct {
	NextRule Rule
}

// Check if a Rule can be applied. This is a Generic Method that can be overridden if necessary.
//
// @param itemConfirm confirmer_interface.ItemConfirm - confirm helper
// @param stocks *map[int]int - map[whId: qty, ...] stock values to confirm from for current item
// @return bool
//
func (r AbstractRule) CanApply(itemConfirm confirmer_interface.ItemConfirm, stocks *map[int]int) bool {
	if itemConfirm.GetConfirmedQty() == itemConfirm.GetNeededQty() {
		return false
	}

	if stocks == nil || 0 == len(*stocks) {
		return false
	}

	return true
}

// This method should be implemented in each specific Rule.
//
// @param itemConfirm confirmer_interface.ItemConfirm - confirm helper
// @param firstGroup *map[int]int - map[whId: qty, ...] stock values to confirm from for current item
// @param secondGroup *map[int]int - map[whId: qty, ...] stock values to confirm from for current item
// @return Rule
//
func (r AbstractRule) Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	panic("This is abstract method - please provide implementation")
	return r
}

// Each specific Rule decides when to apply the next rule.
//
// @param itemConfirm confirmer_interface.ItemConfirm - confirm helper
// @param firstGroup *map[int]int - map[whId: qty, ...] stock values for current item
// @param secondGroup *map[int]int - map[whId: qty, ...] stock values for current item
// @return Rule
//
func (r AbstractRule) ApplyNextRule(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule {
	if nil != r.NextRule {
		r.NextRule.Apply(itemConfirm, firstGroup, secondGroup)
	}

	return r
}
