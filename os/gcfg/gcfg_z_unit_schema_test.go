// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"
)

// testServerConfig is a simplified version of ghttp.ServerConfig for testing.
type testServerConfig struct {
	Name         string        `json:"name"         d:"default"   v:"required"  dc:"Server name|i18n:config.server.name"`
	Address      string        `json:"address"      d:":0"        v:"required"  dc:"Server listening address|i18n:config.server.address"`
	ReadTimeout  time.Duration `json:"readTimeout"  d:"60s"                     dc:"HTTP read timeout|i18n:config.server.readTimeout"`
	KeepAlive    bool          `json:"keepAlive"    d:"true"                    dc:"Enable HTTP keep-alive"`
	unexported   string        // should be skipped
}

// TestBaseConfig tests embedded struct scanning.
type TestBaseConfig struct {
	Host string `json:"host" d:"localhost" dc:"Hostname|i18n:config.base.host"`
	Port int    `json:"port" d:"3306"      dc:"Port number"`
}

type TestDatabaseConfig struct {
	TestBaseConfig          // embedded
	User           string   `json:"user"     d:"root"  v:"required" dc:"Database user|i18n:config.database.user"`
	Password       string   `json:"password"           v:"required" dc:"Database password|i18n:config.database.password"`
}

func TestSchemaRegistry_Register(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		groupMap := map[string]string{
			"Name":        "Basic",
			"Address":     "Basic",
			"ReadTimeout": "Timeout",
			"KeepAlive":   "Basic",
		}

		registry.Register("server", "server", testServerConfig{}, groupMap)

		schema, ok := registry.Get("server")
		t.Assert(ok, true)
		t.AssertNE(schema, nil)
		t.Assert(schema.Name, "server")
		t.Assert(schema.ConfigNode, "server")
		t.Assert(len(schema.Fields) > 0, true)

		// Check groups are extracted correctly.
		t.Assert(len(schema.Groups), 2) // Basic, Timeout
		t.Assert(schema.Groups[0], "Basic")
		t.Assert(schema.Groups[1], "Timeout")
	})
}

func TestSchemaRegistry_FieldParsing(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		groupMap := map[string]string{
			"Name":        "Basic",
			"Address":     "Basic",
			"ReadTimeout": "Timeout",
			"KeepAlive":   "Basic",
		}

		registry.Register("server", "server", testServerConfig{}, groupMap)

		schema, _ := registry.Get("server")

		// Find the Name field.
		var nameField *gcfg.FieldSchema
		for _, f := range schema.Fields {
			if f.Name == "Name" {
				nameField = f
				break
			}
		}
		t.AssertNE(nameField, nil)
		t.Assert(nameField.JsonKey, "name")
		t.Assert(nameField.Type, "string")
		t.Assert(nameField.Default, "default")
		t.Assert(nameField.Rule, "required")
		t.Assert(nameField.Description, "Server name")
		t.Assert(nameField.I18nKey, "config.server.name")
		t.Assert(nameField.Group, "Basic")

		// Find the ReadTimeout field (duration type).
		var timeoutField *gcfg.FieldSchema
		for _, f := range schema.Fields {
			if f.Name == "ReadTimeout" {
				timeoutField = f
				break
			}
		}
		t.AssertNE(timeoutField, nil)
		t.Assert(timeoutField.Type, "duration")
		t.Assert(timeoutField.Default, "60s")
		t.Assert(timeoutField.Group, "Timeout")

		// Find the KeepAlive field (bool type).
		var keepAliveField *gcfg.FieldSchema
		for _, f := range schema.Fields {
			if f.Name == "KeepAlive" {
				keepAliveField = f
				break
			}
		}
		t.AssertNE(keepAliveField, nil)
		t.Assert(keepAliveField.Type, "bool")
		t.Assert(keepAliveField.Default, "true")

		// Unexported field should NOT be present.
		for _, f := range schema.Fields {
			t.AssertNE(f.Name, "unexported")
		}
	})
}

func TestSchemaRegistry_EmbeddedStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		groupMap := map[string]string{
			"Host":     "Connection",
			"Port":     "Connection",
			"User":     "Auth",
			"Password": "Auth",
		}

		registry.Register("database", "database", TestDatabaseConfig{}, groupMap)

		schema, ok := registry.Get("database")
		t.Assert(ok, true)

		// Should have 4 fields: Host, Port from embedded + User, Password from own fields.
		t.Assert(len(schema.Fields), 4)

		// Check Host field from embedded struct.
		var hostField *gcfg.FieldSchema
		for _, f := range schema.Fields {
			if f.Name == "Host" {
				hostField = f
				break
			}
		}
		t.AssertNE(hostField, nil)
		t.Assert(hostField.Default, "localhost")
		t.Assert(hostField.I18nKey, "config.base.host")
		t.Assert(hostField.Group, "Connection")
	})
}

func TestSchemaRegistry_GetAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		registry.Register("server", "server", testServerConfig{}, nil)
		registry.Register("database", "database", TestDatabaseConfig{}, nil)

		all := registry.GetAll()
		t.Assert(len(all), 2)

		// Registration order is maintained.
		t.Assert(all[0].Name, "server")
		t.Assert(all[1].Name, "database")
	})
}

func TestSchemaRegistry_GetNonExistent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		schema, ok := registry.Get("nonexistent")
		t.Assert(ok, false)
		t.Assert(schema, nil)
	})
}

func TestSchemaRegistry_GlobalRegistry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test global functions.
		gcfg.RegisterSchema("test_module", "test", testServerConfig{}, map[string]string{
			"Name": "Basic",
		})

		schema, ok := gcfg.GetSchema("test_module")
		t.Assert(ok, true)
		t.AssertNE(schema, nil)
		t.Assert(schema.Name, "test_module")

		all := gcfg.GetAllSchemas()
		t.Assert(len(all) > 0, true)

		// Global registry should be accessible.
		reg := gcfg.GetGlobalRegistry()
		t.AssertNE(reg, nil)
	})
}

func TestSchemaRegistry_DcTagParsing(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		type testConfig struct {
			Field1 string `json:"field1" dc:"Some description|i18n:config.test.field1"`
			Field2 string `json:"field2" dc:"Just a description"`
			Field3 string `json:"field3"`
		}

		registry.Register("test", "test", testConfig{}, nil)
		schema, _ := registry.Get("test")

		// Field1: description + i18n.
		t.Assert(schema.Fields[0].Description, "Some description")
		t.Assert(schema.Fields[0].I18nKey, "config.test.field1")

		// Field2: description only.
		t.Assert(schema.Fields[1].Description, "Just a description")
		t.Assert(schema.Fields[1].I18nKey, "")

		// Field3: empty.
		t.Assert(schema.Fields[2].Description, "")
		t.Assert(schema.Fields[2].I18nKey, "")
	})
}

func TestSchemaRegistry_PointerStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		// Register with pointer to struct.
		registry.Register("server_ptr", "server", &testServerConfig{}, nil)

		schema, ok := registry.Get("server_ptr")
		t.Assert(ok, true)
		t.Assert(len(schema.Fields) > 0, true)
	})
}

func TestSchemaRegistry_MapType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		type testConfig struct {
			Data map[string]any `json:"data" dc:"Map data"`
			Tags []string       `json:"tags" dc:"Tag list"`
		}

		registry.Register("maptest", "test", testConfig{}, nil)
		schema, _ := registry.Get("maptest")

		// Map type.
		var dataField *gcfg.FieldSchema
		for _, f := range schema.Fields {
			if f.Name == "Data" {
				dataField = f
				break
			}
		}
		t.AssertNE(dataField, nil)
		t.Assert(dataField.Type, "map")

		// Slice type.
		var tagsField *gcfg.FieldSchema
		for _, f := range schema.Fields {
			if f.Name == "Tags" {
				tagsField = f
				break
			}
		}
		t.AssertNE(tagsField, nil)
		t.Assert(tagsField.Type, "[]string")
	})
}

func TestSchemaRegistry_NilGroupMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		registry := gcfg.NewSchemaRegistry()

		// Register with nil groupMap — all fields should be "Other".
		registry.Register("nogroup", "test", testServerConfig{}, nil)
		schema, _ := registry.Get("nogroup")

		for _, f := range schema.Fields {
			t.Assert(f.Group, "Other")
		}
	})
}
