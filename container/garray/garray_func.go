// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"sort"
	"strings"
)

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

func defaultSorter[T comparable](values []T, comparator func(a T, b T) int) {
	sort.Slice(values, func(i, j int) bool {
		return comparator(values[i], values[j]) < 0
	})
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

// tToAnySlice converts []T to []any
func tToAnySlice[T any](values []T) []any {
	if values == nil {
		return nil
	}
	anyValues := make([]any, len(values), cap(values))
	for k, v := range values {
		anyValues[k] = v
	}
	return anyValues
}

// anyToTSlice is convert []any to []T
func anyToTSlice[T any](values []any) []T {
	if values == nil {
		return nil
	}
	tValues := make([]T, len(values), cap(values))
	for k, v := range values {
		tValues[k], _ = v.(T)
	}
	return tValues
}

// tToAnySlices converts [][]T to [][]any
func tToAnySlices[T any](values [][]T) [][]any {
	if values == nil {
		return nil
	}
	anyValues := make([][]any, len(values), cap(values))
	for k, v := range values {
		anyValues[k] = tToAnySlice(v)
	}
	return anyValues
}

// anyToTSlices converts [][]any to [][]T
func anyToTSlices[T any](values [][]any) [][]T {
	if values == nil {
		return nil
	}
	tValues := make([][]T, len(values), cap(values))
	for k, v := range values {
		tValues[k] = anyToTSlice[T](v)
	}
	return tValues
}
