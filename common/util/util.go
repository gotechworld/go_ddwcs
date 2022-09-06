package util

import (
	"sort"
)

/**
 * Implode a map - the final string will contain both key and value
 */
func Implode(mapToImplode map[string]string) string {
	if len(mapToImplode) == 0 {
		return ""
	}

	var keys []string
	for k, _ := range mapToImplode {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	str := ""
	for key, value := range mapToImplode {
		str = str + "_" + key + "_" + value
	}

	return str
}

func ArrayIntersect(array1, array2 []int) []int {
	len1 := len(array1)
	len2 := len(array2)

	if len1 == 0 || len2 == 0 {
		return []int{}
	}

	first := array1
	second := array2

	// check for the longer element
	if len1 < len2 {
		first = array2
		second = array1
	}

	intersection := make([]int, 0, len(first))

	for _, item1 := range first {
		for _, item2 := range second {
			if item1 == item2 {
				intersection = append(intersection, item1)
				break
			}
		}
	}

	return intersection
}

func ArrayIntersectMultiple(arrays ...[]int) []int {
	if len(arrays) == 0 {
		return []int{}
	}

	if len(arrays) == 1 {
		return arrays[0]
	}

	//intersect first two elements
	intersection := ArrayIntersect(arrays[0], arrays[1])

	if len(intersection) == 0 || len(arrays) == 2 {
		return intersection
	}

	for _, item := range arrays[2:] {
		intersection = ArrayIntersect(intersection, item)

		if len(intersection) == 0 {
			break
		}
	}

	return intersection
}
