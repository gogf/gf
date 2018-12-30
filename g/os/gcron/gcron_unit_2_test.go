// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcron_test

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gcron"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

var (
    singletonCron = gcron.New()
)

func TestCron_AddSingleton(t *testing.T) {
    array := garray.New(0, 0)
    singletonCron.AddSingleton("* * * * * *", func() {
        array.Append(1)
        time.Sleep(5*time.Second)

    })
    gtest.Assert(len(singletonCron.Entries()), 1)
    time.Sleep(3500*time.Millisecond)
    gtest.Assert(array.Len(), 1)
}
