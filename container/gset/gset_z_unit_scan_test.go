// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gset_test

import (
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func TestArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		set := gset.New(false)
		set.Add("1")
		set.Add("2")
		set.Add("3")
		array := gset.ScanGSet[string](set)
		t.AssertIN(array[0], []string{"1", "2", "3"})
		t.AssertIN(array[1], []string{"1", "2", "3"})
		t.AssertIN(array[3], []string{"1", "2", "3"})
	})
	gtest.C(t, func(t *gtest.T) {
		set := gset.New(false)
		set.Add(1)
		set.Add(2)
		set.Add(3)
		array := gset.ScanGSet[int](set)
		t.AssertIN(array[0], []int{1, 2, 3})
		t.AssertIN(array[1], []int{1, 2, 3})
		t.AssertIN(array[2], []int{1, 2, 3})
	})

}
