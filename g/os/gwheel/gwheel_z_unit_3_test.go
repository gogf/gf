// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 指定次数运行测试

package gwheel_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gwheel"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestWheel_Times(t *testing.T) {
    wheel := gwheel.NewDefault()
    array := garray.New(0, 0)
    entry := wheel.AddTimes(10, 20, func() {
        array.Append(1)
    })
    gtest.AssertNE(entry, nil)
    time.Sleep(3500*time.Millisecond)
    gtest.Assert(array.Len(), 2)
}
