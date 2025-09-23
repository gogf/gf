// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func TestWatcher_AddWatcherAndNotify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			m          = make(map[string]bool)
			key1       = "test-watcher1"
			key2       = "test-watcher2"
			configFile = guid.S() + ".toml"
			content1   = `key = "value1"`
			content2   = `key = "value2"`
		)

		// Create config file.
		err := gfile.PutContents(configFile, content1)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create config instance.
		c, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)
		m[key1] = true
		m[key2] = true

		// Add watchers.
		c.AddWatcher(key1, func() {
			m[key1] = false
		})
		c.AddWatcher(key2, func() {
			m[key2] = false
		})

		// Check initial values.
		t.Assert(c.MustGet(ctx, "key").String(), "value1")
		t.Assert(m[key1], true)
		t.Assert(m[key2], true)

		// Update config file content.
		err = gfile.PutContents(configFile, content2)
		t.AssertNil(err)

		// Wait for watching notification.
		time.Sleep(1 * time.Second)

		// Check updated values.
		t.Assert(c.MustGet(ctx, "key").String(), "value2")
		t.AssertEQ(m[key1], false)
		t.AssertEQ(m[key2], false)
	})
}

func TestWatcher_RemoveWatcher(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			m          = make(map[string]bool)
			key1       = "test-watcher1"
			key2       = "test-watcher2"
			configFile = guid.S() + ".toml"
			content1   = `key = "value1"`
			content2   = `key = "value2"`
		)
		err := gfile.PutContents(configFile, content1)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create config instance.
		c, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)
		m[key1] = true
		m[key2] = true

		// Add watchers.
		c.AddWatcher(key1, func() {
			m[key1] = false
		})
		c.AddWatcher(key2, func() {
			m[key2] = false
		})

		// Check initial values.
		t.Assert(c.MustGet(ctx, "key").String(), "value1")
		t.Assert(m[key1], true)
		t.Assert(m[key2], true)

		// Remove one watcher.
		c.RemoveWatcher(key2)

		// Update config file content.
		err = gfile.PutContents(configFile, content2)
		t.AssertNil(err)

		// Wait for watching notification.
		time.Sleep(1 * time.Second)

		// Check updated values.
		t.Assert(c.MustGet(ctx, "key").String(), "value2")
		t.AssertEQ(m[key1], false)
		// watcherName2 should not be notified as it was removed
		t.AssertEQ(m[key2], true)
	})
}

func TestWatcher_SetContentNotify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			count      = gvar.New(0, true)
			key        = "test-watcher"
			configFile = guid.S() + ".toml"
			content1   = `key = "value1"`
			content2   = `key = "value2"`
		)

		// Create config file.
		err := gfile.PutContents(configFile, content1)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create config instance.
		c, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)

		// Add watcher.
		c.AddWatcher(key, func() {
			count.Set(count.Int() + 1)
		})

		// Check initial values.
		t.Assert(c.MustGet(ctx, "key").String(), "value1")
		t.Assert(count.Int(), 0)

		// Set custom content.
		c.SetContent(content2, configFile)

		// Wait for watching notification.
		time.Sleep(2 * time.Second)

		// Check that watcher was notified
		t.Assert(count.Int(), 1)
		t.Assert(c.MustGet(ctx, "key").String(), "value2")
	})
}

func TestWatcher_RemoveContentNotify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			count      = gvar.New(0, true)
			key        = "test-watcher"
			configFile = guid.S() + ".toml"
			content    = `key = "value1"`
		)

		// Create config file.
		err := gfile.PutContents(configFile, content)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create config instance.
		c, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)

		// Add watcher.
		c.AddWatcher(key, func() {
			count.Set(count.Int() + 1)
		})

		// Check initial values.
		t.Assert(c.MustGet(ctx, "key").String(), "value1")
		t.Assert(count.Int(), 0)

		// Remove custom content.
		c.RemoveContent(configFile)

		// Wait for watching notification.
		time.Sleep(1 * time.Second)

		// Check that watcher was notified again
		t.Assert(count.Int(), 1)
		t.Assert(c.MustGet(ctx, "key").String(), "value1") // Back to file content
	})
}

func TestWatcher_ClearContentNotify(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			count      = gvar.New(0)
			key        = "test-watcher"
			configFile = guid.S() + ".toml"
			content    = `key = "value1"`
		)

		// Create config file.
		err := gfile.PutContents(configFile, content)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create config instance.
		c, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)

		// Add watcher.
		c.AddWatcher(key, func() {
			count.Set(count.Int() + 1)
		})

		// Check initial values.
		t.Assert(c.MustGet(ctx, "key").String(), "value1")
		t.Assert(count.Int(), 0)

		// Clear all custom content.
		c.ClearContent()

		// Wait for watching notification.
		time.Sleep(1 * time.Second)

		// Check that watcher was notified again
		t.Assert(count.Int(), 1)
		t.Assert(c.MustGet(ctx, "key").String(), "value1") // Back to file content
	})
}
