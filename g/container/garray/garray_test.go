// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go

package garray_test

import (
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
    "testing"
)


func TestArray_Unique(t *testing.T) {
    expect := []int{1, 2, 3, 4, 5, 6}
    array  := garray.NewIntArray(0, 0)
    array.Append(1, 1, 2, 3, 3, 4, 4, 5, 5, 6, 6)
    array.Unique()
    if fmt.Sprint(array.Slice()) != fmt.Sprint(expect) {
        t.Errorf("get: %v, expect: %v\n", array.Slice(), expect)
    }
}
