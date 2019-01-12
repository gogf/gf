// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 指定次数运行测试

package gtimer_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestTimer_Times(t *testing.T) {
    gtest.Case(t, func() {
        wheel := New()
        array := garray.New(0, 0)
        wheel.AddTimes(time.Second, 2, func() {
            array.Append(1)
        })
        time.Sleep(3500*time.Millisecond)
        gtest.Assert(array.Len(), 2)
    })
}
