// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Core_SetDebug_GetDebug(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := configs.config
		defer func() {
			configs.config = originalConfig
		}()

		// Create a test configuration
		configs.config = make(Config)
		testNode := ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "test_db",
			Type: "mysql",
		}
		err := AddConfigNode("test_group", testNode)
		t.AssertNil(err)

		// Create Core instance
		node, err := GetConfigGroup("test_group")
		t.AssertNil(err)
		core := &Core{
			group:  "test_group",
			config: &node[0],
			debug:  gtype.NewBool(false),
		}

		// Test default value
		result := core.GetDebug()
		t.Assert(result, false)

		// Test setting debug to true
		core.SetDebug(true)
		result = core.GetDebug()
		t.Assert(result, true)

		// Test setting debug to false
		core.SetDebug(false)
		result = core.GetDebug()
		t.Assert(result, false)
	})
}

func Test_Core_SetDryRun_GetDryRun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := configs.config
		defer func() {
			configs.config = originalConfig
		}()

		// Create a test configuration
		configs.config = make(Config)
		testNode := ConfigNode{
			Host:   "127.0.0.1",
			Port:   "3306",
			User:   "root",
			Pass:   "123456",
			Name:   "test_db",
			Type:   "mysql",
			DryRun: false,
		}
		err := AddConfigNode("test_group", testNode)
		t.AssertNil(err)

		// Create Core instance
		node, err := GetConfigGroup("test_group")
		t.AssertNil(err)
		core := &Core{
			group:  "test_group",
			config: &node[0],
		}

		// Test default value
		result := core.GetDryRun()
		t.Assert(result, false)

		// Test setting dry run to true
		core.SetDryRun(true)
		result = core.GetDryRun()
		t.Assert(result, true)

		// Test setting dry run to false
		core.SetDryRun(false)
		result = core.GetDryRun()
		t.Assert(result, false)
	})
}

func Test_Core_SetLogger_GetLogger(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create Core instance
		core := &Core{}

		// Test setting custom logger
		customLogger := glog.New()
		core.SetLogger(customLogger)
		result := core.GetLogger()
		t.Assert(result, customLogger)

		// Test setting nil logger
		core.SetLogger(nil)
		result = core.GetLogger()
		t.Assert(result, nil)
	})
}

func Test_Core_SetMaxConnections(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create Core instance
		core := &Core{}

		// Test SetMaxIdleConnCount
		core.SetMaxIdleConnCount(10)
		t.Assert(core.dynamicConfig.MaxIdleConnCount, 10)

		// Test SetMaxOpenConnCount
		core.SetMaxOpenConnCount(20)
		t.Assert(core.dynamicConfig.MaxOpenConnCount, 20)

		// Test SetMaxConnLifeTime
		testDuration := time.Hour
		core.SetMaxConnLifeTime(testDuration)
		t.Assert(core.dynamicConfig.MaxConnLifeTime, testDuration)
	})
}

func Test_Core_GetCache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create Core instance
		core := &Core{}

		cache := core.GetCache()
		// Cache might be nil if not initialized, so we just test that the call doesn't panic
		_ = cache
	})
}

func Test_Core_GetGroup(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create Core instance
		core := &Core{
			group: "test_group",
		}

		group := core.GetGroup()
		t.Assert(group, "test_group")
	})
}

func Test_Core_GetPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := configs.config
		defer func() {
			configs.config = originalConfig
		}()

		// Create a test configuration
		configs.config = make(Config)
		testNode := ConfigNode{
			Host:   "127.0.0.1",
			Port:   "3306",
			User:   "root",
			Pass:   "123456",
			Name:   "test_db",
			Type:   "mysql",
			Prefix: "gf_",
		}
		err := AddConfigNode("test_group", testNode)
		t.AssertNil(err)

		// Create Core instance
		node, err := GetConfigGroup("test_group")
		t.AssertNil(err)
		core := &Core{
			group:  "test_group",
			config: &node[0],
		}

		prefix := core.GetPrefix()
		t.Assert(prefix, "gf_")
	})
}
