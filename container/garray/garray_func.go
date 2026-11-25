// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"slices"
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

// defaultSorter is a generic sorting function that sorts a slice of comparable types
// using the provided comparator function.
func defaultSorter[T comparable](values []T, comparator func(a T, b T) int) {
	slices.SortFunc(values, comparator)
}

// cloneSlice returns a non-nil copy of the given slice.
func cloneSlice[T any](values []T) []T {
	cloned := slices.Clone(values)
	if cloned == nil {
		cloned = make([]T, 0)
	}
	return cloned
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
