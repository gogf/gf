// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_OTEL_SQLTracing_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "sqlite",
			Name: ":memory:",
		}
		
		// By default, SQL tracing should be disabled
		t.Assert(config.OtelTraceSQLEnabled, false)
	})
}

func Test_OTEL_SQLTracing_Configuration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "sqlite",
			Name: ":memory:",
			OtelTraceSQLEnabled: true,
		}
		
		// SQL tracing should be configurable
		t.Assert(config.OtelTraceSQLEnabled, true)
	})
}

func Test_OTEL_SQLTracing_Enabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "mysql",
			Name: "test_db",
			OtelTraceSQLEnabled: true,
		}
		
		// Test that the configuration field can be set and retrieved
		t.Assert(config.OtelTraceSQLEnabled, true)
		
		// Test that the field is preserved during configuration operations
		configGroup := gdb.ConfigGroup{config}
		t.Assert(configGroup[0].OtelTraceSQLEnabled, true)
	})
}

func Test_OTEL_SQLTracing_Disabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "mysql", 
			Name: "test_db",
			OtelTraceSQLEnabled: false,
		}
		
		// Test that the configuration field can be set and retrieved
		t.Assert(config.OtelTraceSQLEnabled, false)
		
		// Test that the field is preserved during configuration operations
		configGroup := gdb.ConfigGroup{config}
		t.Assert(configGroup[0].OtelTraceSQLEnabled, false)
	})
}