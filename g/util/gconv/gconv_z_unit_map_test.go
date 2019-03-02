// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/test/gtest"
    "github.com/gogf/gf/g/util/gconv"
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
        gtest.Assert(gconv.Map(m1), g.Map{
            "k" : "v",
        })
        gtest.Assert(gconv.Map(m2), g.Map{
            "3" : "v",
        })
        gtest.Assert(gconv.Map(m3), g.Map{
            "1.22" : "3.1",
        })
    })
}

