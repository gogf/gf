// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestCron_Entry_Operations(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			cron  = gcron.New()
			array = garray.New(true)
		)
		go func() {
			cron.DelayAddTimes(ctx, 500*time.Millisecond, "* * * * * *", 2, func(ctx context.Context) {
				array.Append(1)
			})
			t.Assert(cron.Size(), 0)
		}()
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
		var entry *gcron.Entry
		go func() {
			var err error
			entry, err = cron.Add(ctx, "* * * * * *", func(ctx context.Context) {
				array.Append(1)
			})
			t.Assert(err, nil)
			t.Assert(array.Len(), 0)
			t.Assert(cron.Size(), 1)
		}()
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Stop()
		time.Sleep(5000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Start()
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 1)
		entry.Close()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 0)
	})
}
