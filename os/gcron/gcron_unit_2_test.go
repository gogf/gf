// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron_test

import (
	"github.com/gogf/gf/frame/g"
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/test/gtest"
)

func TestCron_Job_Operations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			cron  = gcron.New()
			array = garray.New(true)
		)
		cron.DelayAddTimes(500*time.Millisecond, "* * * * * *", 2, func() {
			g.Log().Println("add times")
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(3000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			cron  = gcron.New()
			array = garray.New(true)
		)
		entry, err1 := cron.Add("* * * * * *", func() {
			g.Log().Println("add")
			array.Append(1)
		})
		t.Assert(err1, nil)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Stop()
		time.Sleep(5000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Start()
		g.Log().Println("start")
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 1)
		entry.Close()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 0)
	})
}
