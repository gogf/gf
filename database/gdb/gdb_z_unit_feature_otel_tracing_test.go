// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/internal/otel"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_OTEL_SQLTracing_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "sqlite",
			Name: ":memory:",
		}
		
		// By default, SQL tracing should be disabled
		t.Assert(config.IsOtelTraceSQLEnabled(), false)
		t.Assert(config.OtelTraceSQLEnabled, false)
		t.Assert(config.Otel.TraceSQLEnabled, false)
	})
}

func Test_OTEL_SQLTracing_Configuration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "sqlite",
			Name: ":memory:",
			OtelTraceSQLEnabled: true,
		}
		
		// SQL tracing should be configurable using legacy field
		t.Assert(config.IsOtelTraceSQLEnabled(), true)
		t.Assert(config.OtelTraceSQLEnabled, true)
		t.Assert(config.Otel.TraceSQLEnabled, false)
	})
}

func Test_OTEL_SQLTracing_NewConfiguration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "sqlite",
			Name: ":memory:",
			Otel: otel.Config{
				TraceSQLEnabled: true,
			},
		}
		
		// SQL tracing should be configurable using new configuration
		t.Assert(config.IsOtelTraceSQLEnabled(), true)
		t.Assert(config.OtelTraceSQLEnabled, false)
		t.Assert(config.Otel.TraceSQLEnabled, true)
	})
}

func Test_OTEL_SQLTracing_Enabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "mysql",
			Name: "test_db",
			OtelTraceSQLEnabled: true,
		}
		
		// Test that the configuration field can be set and retrieved using legacy field
		t.Assert(config.IsOtelTraceSQLEnabled(), true)
	})
}

func Test_OTEL_SQLTracing_BothFieldsEnabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := gdb.ConfigNode{
			Type: "mysql",
			Name: "test_db",
			OtelTraceSQLEnabled: false,
			Otel: otel.Config{
				TraceSQLEnabled: true,
			},
		}
		
		// New field should take precedence over legacy field
		t.Assert(config.IsOtelTraceSQLEnabled(), true)
	})
}