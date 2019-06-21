// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron_test

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/gcron"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

func TestCron_Entry_Operations(t *testing.T) {
	gtest.Case(t, func() {

		gtest.Case(t, func() {
			cron := gcron.New()
			array := garray.New()
			cron.DelayAddTimes(500*time.Millisecond, "* * * * * *", 2, func() {
				glog.Println("add times")
				array.Append(1)
			})
			gtest.Assert(cron.Size(), 0)
			time.Sleep(800 * time.Millisecond)
			gtest.Assert(array.Len(), 0)
			gtest.Assert(cron.Size(), 1)
			time.Sleep(3000 * time.Millisecond)
			gtest.Assert(array.Len(), 2)
			gtest.Assert(cron.Size(), 0)
		})

		cron := gcron.New()
		array := garray.New()
		entry, err1 := cron.Add("* * * * * *", func() {
			glog.Println("add")
			array.Append(1)
		})
		gtest.Assert(err1, nil)
		gtest.Assert(array.Len(), 0)
		gtest.Assert(cron.Size(), 1)
		time.Sleep(1200 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		gtest.Assert(cron.Size(), 1)
		entry.Stop()
		time.Sleep(2000 * time.Millisecond)
		gtest.Assert(array.Len(), 1)
		gtest.Assert(cron.Size(), 1)
		entry.Start()
		glog.Println("start")
		time.Sleep(1200 * time.Millisecond)
		gtest.Assert(array.Len(), 2)
		gtest.Assert(cron.Size(), 1)
		entry.Close()
		time.Sleep(1200 * time.Millisecond)
		gtest.Assert(cron.Size(), 0)
	})
}
