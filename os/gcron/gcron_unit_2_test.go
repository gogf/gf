// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gcron_test

import (
	"testing"
	"time"

	"github.com/jin502437344/gf/container/garray"
	"github.com/jin502437344/gf/os/gcron"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/test/gtest"
)

func TestCron_Entry_Operations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cron := gcron.New()
		array := garray.New(true)
		cron.DelayAddTimes(500*time.Millisecond, "* * * * * *", 2, func() {
			glog.Println("add times")
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
		cron := gcron.New()
		array := garray.New(true)
		entry, err1 := cron.Add("* * * * * *", func() {
			glog.Println("add")
			array.Append(1)
		})
		t.Assert(err1, nil)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Stop()
		time.Sleep(2000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Start()
		glog.Println("start")
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 1)
		entry.Close()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 0)
	})
}
