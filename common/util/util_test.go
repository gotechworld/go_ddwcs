package util_test

import (
	"gitlab.altex.ro/ams/go_ddwcs/common/util"
	"testing"
)

func TestArrayIntersect(t *testing.T) {
	//given
	firstSlice := []int{1, 2, 3, 4}
	secondSlice := []int{3, 4, 5, 6}

	//when
	if intersection := util.ArrayIntersect(firstSlice, secondSlice); len(intersection) == 2 {
		t.Log("Passed")
	}

	//then
	t.Log("No common element")
	t.FailNow()
}

func TestArrayIntersectMultiple(t *testing.T) {
	//given
	firstSlice := []int{1, 2, 3, 4}
	secondSlice := []int{3, 4, 5, 6}
	thirdSlice := []int{4, 5, 6, 7}

	//when
	if intersection := util.ArrayIntersectMultiple(firstSlice, secondSlice, thirdSlice); len(intersection) == 1 {
		t.Log("Passed")
	}

	//then
	t.Log("No common element")
	t.FailNow()
}
