// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import "strings"

// defaultComparatorInt for int comparison.
func defaultComparatorInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// defaultComparatorStr for string comparison.
func defaultComparatorStr(a, b string) int {
	return strings.Compare(a, b)
}

// quickSortInt is the quick-sorting algorithm implements for int.
func quickSortInt(values []int, comparator func(a, b int) int) {
	if len(values) <= 1 {
		return
	}
	mid, i := values[0], 1
	head, tail := 0, len(values)-1
	for head < tail {
		if comparator(values[i], mid) > 0 {
			values[i], values[tail] = values[tail], values[i]
			tail--
		} else {
			values[i], values[head] = values[head], values[i]
			head++
			i++
		}
	}
	values[head] = mid
	quickSortInt(values[:head], comparator)
	quickSortInt(values[head+1:], comparator)
}

// quickSortStr is the quick-sorting algorithm implements for string.
func quickSortStr(values []string, comparator func(a, b string) int) {
	if len(values) <= 1 {
		return
	}
	mid, i := values[0], 1
	head, tail := 0, len(values)-1
	for head < tail {
		if comparator(values[i], mid) > 0 {
			values[i], values[tail] = values[tail], values[i]
			tail--
		} else {
			values[i], values[head] = values[head], values[i]
			head++
			i++
		}
	}
	values[head] = mid
	quickSortStr(values[:head], comparator)
	quickSortStr(values[head+1:], comparator)
}

// ScanGArray Convert *garray.Array to a Slice of the specified type
func ScanGArray[T any](arr *Array) []T {
	dst := make([]T, 0)
	arr.Iterator(func(k int, v interface{}) bool {
		dst = append(dst, v.(T))
		return true
	})
	return dst
}

// ScanGSortArray Convert *garray.SortedArray to a Slice of the specified type
func ScanGSortArray[T any](arr *SortedArray) []T {
	dst := make([]T, 0)
	arr.IteratorAsc(func(k int, v interface{}) bool {
		dst = append(dst, v.(T))
		return true
	})
	return dst
}
