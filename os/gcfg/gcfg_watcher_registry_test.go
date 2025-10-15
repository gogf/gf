// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestWatcherRegistry_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewWatcherRegistry()

		// Test Add and GetNames
		called := false
		registry.Add("test-watcher", func(ctx context.Context) {
			called = true
		})

		names := registry.GetNames()
		t.AssertEQ(len(names), 1)
		t.AssertEQ(names[0], "test-watcher")

		// Test Notify
		registry.Notify(context.Background())
		time.Sleep(100 * time.Millisecond)
		t.AssertEQ(called, true)

		// Test Remove
		registry.Remove("test-watcher")
		names = registry.GetNames()
		t.AssertEQ(len(names), 0)
	})
}

func TestWatcherRegistry_MultipleWatchers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewWatcherRegistry()

		var count1, count2, count3 int
		registry.Add("watcher1", func(ctx context.Context) {
			count1++
		})
		registry.Add("watcher2", func(ctx context.Context) {
			count2++
		})
		registry.Add("watcher3", func(ctx context.Context) {
			count3++
		})

		names := registry.GetNames()
		t.AssertEQ(len(names), 3)

		registry.Notify(context.Background())
		time.Sleep(100 * time.Millisecond)
		t.AssertEQ(count1, 1)
		t.AssertEQ(count2, 1)
		t.AssertEQ(count3, 1)

		// Remove one watcher
		registry.Remove("watcher2")
		names = registry.GetNames()
		t.AssertEQ(len(names), 2)
	})
}
