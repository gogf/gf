// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
	"testing"
)

func TestArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		newArray := garray.NewArray(false)
		newArray.Append("1")
		newArray.Append("2")
		newArray.Append("3")
		array := garray.ScanGArray[string](newArray)
		t.Assert(array, []string{"1", "2", "3"})
	})
	gtest.C(t, func(t *gtest.T) {
		newArray := garray.NewArray(false)
		newArray.Append(1)
		newArray.Append(2)
		newArray.Append(3)
		array := garray.ScanGArray[int](newArray)
		t.Assert(array, []int{1, 2, 3})
	})
	gtest.C(t, func(t *gtest.T) {
		newArray := garray.NewSortedArray(gutil.ComparatorString)
		newArray.Append("1")
		newArray.Append("2")
		newArray.Append("3")
		array := garray.ScanGSortArray[string](newArray)
		t.Assert(array, []string{"1", "2", "3"})
	})
	gtest.C(t, func(t *gtest.T) {
		newArray := garray.NewSortedArray(gutil.ComparatorInt)
		newArray.Append(1)
		newArray.Append(2)
		newArray.Append(3)
		array := garray.ScanGSortArray[int](newArray)
		t.Assert(array, []int{1, 2, 3})
	})
}
