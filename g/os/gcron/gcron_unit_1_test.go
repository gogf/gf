// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gcron_test

import (
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/os/gcron"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
    "time"
)

func TestCron_Add_Close(t *testing.T) {
    gtest.Case(func() {
        cron  := gcron.New()
        array := garray.New(0, 0)
        _, err1 := cron.Add("* * * * * *", func() {
            array.Append(1)
        })
        _, err2 := cron.Add("* * * * * *", func() {
            array.Append(1)
        }, "test")
        _, err3 := cron.Add("* * * * * *", func() {
            array.Append(1)
        }, "test")
        _, err4 := cron.Add("@every 2s", func() {
            array.Append(1)
        })
        gtest.Assert(err1, nil)
        gtest.Assert(err2, nil)
        gtest.AssertNE(err3, nil)
        gtest.Assert(err4, nil)
        gtest.Assert(cron.Size(), 3)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 2)
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), 5)
        cron.Close()
        time.Sleep(1100*time.Millisecond)
        fixedLength := array.Len()
        time.Sleep(1100*time.Millisecond)
        gtest.Assert(array.Len(), fixedLength)
    })
}

func TestCron_Method(t *testing.T) {
    gtest.Case(func() {
        cron  := gcron.New()
        cron.Add("* * * * * *", func() {}, "add")
        fmt.Println("start", time.Now())
        cron.DelayAdd(1, "* * * * * *", func() {}, "delay_add")
        gtest.Assert(cron.Size(), 1)
        time.Sleep(1200*time.Millisecond)
        gtest.Assert(cron.Size(), 2)

        cron.Remove("delay_add")
        gtest.Assert(cron.Size(), 1)

        entry1 := cron.Search("add")
        entry2 := cron.Search("test-none")
        gtest.AssertNE(entry1, nil)
        gtest.Assert(entry2, nil)
    })
}
