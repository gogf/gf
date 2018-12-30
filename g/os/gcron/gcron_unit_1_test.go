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
    cron1 = gcron.New()
    cron2 = gcron.New()
)
func TestCron_Add_Close(t *testing.T) {
    array := garray.New(0, 0)
    _, err1 := cron1.Add("* * * * * *", func() {
        array.Append(1)
    })
    _, err2 := cron1.Add("* * * * * *", func() {
        array.Append(1)
    }, "test")
    _, err3 := cron1.Add("* * * * * *", func() {
        array.Append(1)
    }, "test")
    _, err4 := cron1.Add("@every 2s", func() {
        array.Append(1)
    })
    gtest.Assert(err1, nil)
    gtest.Assert(err2, nil)
    gtest.AssertNE(err3, nil)
    gtest.Assert(err4, nil)
    time.Sleep(1100*time.Millisecond)
    gtest.Assert(array.Len(), 2)
    time.Sleep(1100*time.Millisecond)
    gtest.Assert(array.Len(), 5)
    cron1.Close()
    time.Sleep(1100*time.Millisecond)
    fixedLength := array.Len()
    time.Sleep(1100*time.Millisecond)
    gtest.Assert(array.Len(), fixedLength)
}

func TestCron_Entries(t *testing.T) {
    entries := cron1.Entries()
    gtest.Assert(len(entries), 3)
}

func TestCron_DelayAdd(t *testing.T) {
    cron2.Add("* * * * * *", func() {}, "add")
    cron2.DelayAdd(1, "* * * * * *", func() {}, "delay_add")
    gtest.Assert(len(cron2.Entries()), 1)
    time.Sleep(1100*time.Millisecond)
    gtest.Assert(len(cron2.Entries()), 2)
}

func TestCron_Remove(t *testing.T) {
    cron2.Remove("delay_add")
    gtest.Assert(len(cron2.Entries()), 1)
}

func TestCron_Search(t *testing.T) {
    entry1 := cron2.Search("add")
    entry2 := cron2.Search("test-none")
    gtest.AssertNE(entry1, nil)
    gtest.Assert(entry2, nil)
}
