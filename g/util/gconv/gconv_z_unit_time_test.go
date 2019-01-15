// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv_test

import (
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)


func Test_Time(t *testing.T) {
    gtest.Case(t, func() {
        t1 := "2011-10-10 01:02:03.456"
        gtest.AssertEQ(gconv.GTime(t1), gtime.NewFromStr(t1))
        gtest.AssertEQ(gconv.Time(t1), gtime.NewFromStr(t1).Time)
        gtest.AssertEQ(gconv.TimeDuration(100), 100*time.Nanosecond)
    })
}
