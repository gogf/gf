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

func TestCron_Entry_Operations(t *testing.T) {
    gtest.Case(t, func() {
        cron  := gcron.New()
        array := garray.New(0, 0)
        entry, err1 := cron.Add("* * * * * *", func() {
            array.Append(1)
        })
        gtest.Assert(err1, nil)
        gtest.Assert(array.Len(), 0)
        gtest.Assert(cron.Size(), 1)
        time.Sleep(1200*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        gtest.Assert(cron.Size(), 1)
        entry.Stop()
        time.Sleep(1200*time.Millisecond)
        gtest.Assert(array.Len(), 1)
        gtest.Assert(cron.Size(), 1)
        entry.Start()
        time.Sleep(1200*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        gtest.Assert(cron.Size(), 1)
        entry.Close()
        time.Sleep(1200*time.Millisecond)
        gtest.Assert(cron.Size(), 0)

    })
}
