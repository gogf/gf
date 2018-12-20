// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go

package garray_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gtest"
    "strings"
    "testing"
)


func Test_IntArray_Unique(t *testing.T) {
    expect := []int{1, 2, 3, 4, 5, 6}
    array  := garray.NewIntArray(0, 0)
    array.Append(1, 1, 2, 3, 3, 4, 4, 5, 5, 6, 6)
    array.Unique()
    gtest.Assert(array.Slice(), expect)
}

func Test_SortedIntArray1(t *testing.T) {
    expect := []int{0,1,2,3,4,5,6,7,8,9,10}
    array  := garray.NewSortedIntArray(0)
    for i := 10; i > -1; i-- {
        array.Add(i)
    }
    gtest.Assert(array.Slice(), expect)
}

func Test_SortedIntArray2(t *testing.T) {
    expect := []int{0,1,2,3,4,5,6,7,8,9,10}
    array  := garray.NewSortedIntArray(0)
    for i := 0; i <= 10; i++ {
        array.Add(i)
    }
    gtest.Assert(array.Slice(), expect)
}

func Test_SortedStringArray1(t *testing.T) {
    expect := []string{"0","1","10","2","3","4","5","6","7","8","9"}
    array  := garray.NewSortedStringArray(0)
    for i := 10; i > -1; i-- {
        array.Add(gconv.String(i))
    }
    gtest.Assert(array.Slice(), expect)
}

func Test_SortedStringArray2(t *testing.T) {
    expect := []string{"0","1","10","2","3","4","5","6","7","8","9"}
    array  := garray.NewSortedStringArray(0)
    for i := 0; i <= 10; i++ {
        array.Add(gconv.String(i))
    }
    gtest.Assert(array.Slice(), expect)
}

func Test_SortedArray1(t *testing.T) {
    expect := []string{"0","1","10","2","3","4","5","6","7","8","9"}
    array  := garray.NewSortedArray(0, func(v1, v2 interface{}) int {
        return strings.Compare(gconv.String(v1), gconv.String(v2))
    })
    for i := 10; i > -1; i-- {
        array.Add(gconv.String(i))
    }
    gtest.Assert(array.Slice(), expect)
}

func Test_SortedArray2(t *testing.T) {
    expect := []string{"0","1","10","2","3","4","5","6","7","8","9"}
    array  := garray.NewSortedArray(0, func(v1, v2 interface{}) int {
        return strings.Compare(gconv.String(v1), gconv.String(v2))
    })
    for i := 0; i <= 10; i++ {
        array.Add(gconv.String(i))
    }
    gtest.Assert(array.Slice(), expect)
}