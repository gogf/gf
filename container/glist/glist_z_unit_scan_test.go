// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glist_test

import (
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func TestList(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		newArray := glist.New(false)
		newArray.PushBack("1")
		newArray.PushBack("2")
		newArray.PushBack("3")
		array := glist.ScanGList[string](newArray)
		t.Assert(array, []string{"1", "2", "3"})
	})
	gtest.C(t, func(t *gtest.T) {
		newArray := glist.New(false)
		newArray.PushBack(1)
		newArray.PushBack(2)
		newArray.PushBack(3)
		array := glist.ScanGList[int](newArray)
		t.Assert(array, []int{1, 2, 3})
	})
}
