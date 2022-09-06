package rule

import (
	"gitlab.altex.ro/ams/go_ddwcs/delivery/service/confirm/confirmer_interface"
)

type Rule interface {
	CanApply(itemConfirm confirmer_interface.ItemConfirm, stocks *map[int]int) bool
	Apply(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule
	ApplyNextRule(itemConfirm confirmer_interface.ItemConfirm, firstGroup *map[int]int, secondGroup *map[int]int) Rule
}
