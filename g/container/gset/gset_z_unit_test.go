// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go

package gset_test

import (
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/test/gtest"
    "testing"
)

func TestIntSet_Basic(t *testing.T) {
    gtest.Case(t, func() {
        s := gset.NewIntSet()
        s.Add(1).Add(1).Add(2)
        s.BatchAdd([]int{3,4})
        gtest.Assert(s.Size(), 3)
        gtest.Assert(s.Contains(4), true)
    })
}
