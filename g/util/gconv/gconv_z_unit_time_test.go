// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
    "time"
)


func Test_Time(t *testing.T) {
    gtest.Case(t, func() {
        t1 := "2011-10-10 01:02:03.456"
        gtest.AssertEQ(gconv.GTime(t1), gtime.NewFromStr(t1))
        gtest.AssertEQ(gconv.Time(t1), gtime.NewFromStr(t1).Time)
        gtest.AssertEQ(gconv.Duration(100), 100*time.Nanosecond)
    })
}
