// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv_test

import (
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
)


func Test_Map(t *testing.T) {
    gtest.Case(t, func() {
        m1 := map[string]string{
            "k" : "v",
        }
        m2 := map[int]string{
            3 : "v",
        }
        m3 := map[float64]float32{
            1.22 : 3.1,
        }
        gtest.Assert(gconv.Map(m1), m1)
        gtest.Assert(gconv.Map(m2), m2)
        gtest.Assert(gconv.Map(m3), m3)
    })
}
