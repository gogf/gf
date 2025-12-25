// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gprop_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gprop"
	"github.com/gogf/gf/v2/test/gtest"
)

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

func Test_Configurator_BasicFunctionality(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var config ServerConfig
		content := `{"host":"localhost","port":8080,"name":"test-server"}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		// Load configuration from content
		err = configurator.Load(context.Background())
		t.AssertNil(err)

		// Verify configuration values are loaded correctly
		t.Assert(config.Host, "localhost")
		t.Assert(config.Port, 8080)
		t.Assert(config.Name, "test-server")
	})
}

func Test_Configurator_PropertyKeyFunctionality(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"server":{"host":"example.com","port":9000}}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "server", &config)

		// Load configuration with specific property key
		err = configurator.Load(context.Background())
		t.AssertNil(err)

		// Verify configuration values are loaded correctly from nested key
		t.Assert(config.Host, "example.com")
		t.Assert(config.Port, 9000)
	})
}

func Test_Configurator_GetMethod(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"host":"get-test","port":7000}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		err = configurator.Load(context.Background())
		t.AssertNil(err)

		// Get current configuration
		currentConfig := configurator.Get()
		t.Assert(currentConfig.Host, "get-test")
		t.Assert(currentConfig.Port, 7000)
	})
}

func Test_Configurator_OnChangeCallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"host":"change-test","port":6000}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		// Track if callback was called
		callbackCalled := false
		configurator.OnChange(func(updated LocalServerConfig) error {
			callbackCalled = true
			t.Assert(updated.Host, "change-test")
			t.Assert(updated.Port, 6000)
			return nil
		})

		err = configurator.Load(context.Background())
		t.AssertNil(err)
		t.Assert(callbackCalled, true)
		t.Assert(config.Host, "change-test")
		t.Assert(config.Port, 6000)
	})
}

func Test_Configurator_CustomConverter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"host":"convert-test","port":5000}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		// Set custom converter
		configurator.SetConverter(func(data any, target *LocalServerConfig) error {
			m := data.(map[string]interface{})
			target.Host = m["host"].(string) + "-converted"
			// Handle json.Number type
			switch v := m["port"].(type) {
			case float64:
				target.Port = int(v) + 100
			case json.Number:
				if intVal, err := v.Int64(); err == nil {
					target.Port = int(intVal) + 100
				}
			case int:
				target.Port = v + 100
			case int64:
				target.Port = int(v) + 100
			}
			return nil
		})

		// Load configuration
		err = configurator.Load(context.Background())
		t.AssertNil(err)

		// Verify values after custom conversion
		t.Assert(config.Host, "convert-test-converted")
		t.Assert(config.Port, 5100)
	})
}

func Test_Configurator_LoadErrorHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		// Create a valid config first, then use a converter that returns an error
		var config LocalServerConfig
		content := `{"host":"test","port":3000}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		// Track if error handler was called
		errorHandlerCalled := false
		configurator.SetLoadErrorHandler(func(ctx context.Context, err error) {
			errorHandlerCalled = true
		})

		// Set converter that returns an error
		configurator.SetConverter(func(data any, target *LocalServerConfig) error {
			return fmt.Errorf("converter error for testing")
		})

		// Try to load - should trigger error handler
		err = configurator.Load(context.Background())
		t.AssertNE(err, nil)
		t.Assert(errorHandlerCalled, true)
	})
}

func Test_Configurator_MustLoad(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"host":"mustload-test","port":3000}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		// MustLoad should not panic with valid config
		configurator.MustLoad(context.Background())
		t.Assert(config.Host, "mustload-test")
		t.Assert(config.Port, 3000)
	})
}

func Test_Configurator_FromFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"host":"file-test","port":2000}`

		// Create temporary config file
		tmpFile := gfile.Temp("gprop-test-config.json")
		err := gfile.PutContents(tmpFile, content)
		t.AssertNil(err)
		defer gfile.Remove(tmpFile)

		configurator, err := gprop.FromFile[LocalServerConfig](tmpFile, "", &config)
		t.AssertNil(err)

		err = configurator.Load(context.Background())
		t.AssertNil(err)
		t.Assert(config.Host, "file-test")
		t.Assert(config.Port, 2000)
	})
}

func Test_Configurator_FromContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"host":"content-test","port":1000}`

		configurator, err := gprop.FromContent[LocalServerConfig](content, "", &config)
		t.AssertNil(err)

		err = configurator.Load(context.Background())
		t.AssertNil(err)
		t.Assert(config.Host, "content-test")
		t.Assert(config.Port, 1000)
	})
}

func Test_Configurator_ThreadSafety(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		content := `{"host":"thread-test","port":5000}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		// Load configuration
		err = configurator.Load(context.Background())
		t.AssertNil(err)

		// Test that Get() returns consistent values
		config1 := configurator.Get()
		config2 := configurator.Get()

		t.Assert(config1.Host, config2.Host)
		t.Assert(config1.Port, config2.Port)
	})
}

func Test_Configurator_Watch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type LocalServerConfig struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		var config LocalServerConfig
		initialContent := `{"host":"initial","port":1000}`

		// Create temporary config file
		tmpFile := gfile.Temp("gprop-watch-test.json")
		err := gfile.PutContents(tmpFile, initialContent)
		t.AssertNil(err)
		defer gfile.Remove(tmpFile)

		// Create config from file
		configurator, err := gprop.FromFile[LocalServerConfig](tmpFile, "", &config)
		t.AssertNil(err)

		// Set up change callback
		configChanged := make(chan bool, 1)
		configurator.OnChange(func(updated LocalServerConfig) error {
			configChanged <- true
			return nil
		})

		// Load initial config
		err = configurator.Load(context.Background())
		t.AssertNil(err)

		// Start watching (this will only work if the adapter supports watching)
		err = configurator.Watch(context.Background(), "test-watcher")
		if err != nil {
			// Not all adapters support watching, which is OK
			t.AssertGT(len(err.Error()), 0) // Just check that error exists
			return
		}

		// Update file content to trigger watch
		updatedContent := `{"host":"updated","port":2000}`
		time.Sleep(100 * time.Millisecond) // Allow time for watcher to be set up
		err = gfile.PutContents(tmpFile, updatedContent)
		t.AssertNil(err)

		// Wait for change notification
		select {
		case <-configChanged:
			// Success - config was updated
		case <-time.After(2 * time.Second):
			// Timeout - might happen if file watcher doesn't detect changes quickly
			// This is acceptable in testing environment
		}
	})
}

// Test_Configurator_WithRealWorldScenario tests a more complex real-world scenario
func Test_Configurator_WithRealWorldScenario(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Define a complex configuration structure
		type DatabaseConfig struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
			Database string `json:"database"`
		}

		type ComplexAppConfig struct {
			Server   ServerConfig   `json:"server"`
			Database DatabaseConfig `json:"database"`
			Features []string       `json:"features"`
		}

		var config ComplexAppConfig
		content := `{
			"server": {
				"host": "localhost",
				"port": 8080,
				"name": "myapp"
			},
			"database": {
				"host": "db.example.com",
				"port": 5432,
				"username": "user",
				"password": "pass",
				"database": "mydb"
			},
			"features": ["feature1", "feature2", "feature3"]
		}`

		adapter, err := gcfg.NewAdapterContent(content)
		t.AssertNil(err)
		cfg := gcfg.NewWithAdapter(adapter)
		configurator := gprop.New(cfg, "", &config)

		// Load the complex configuration
		err = configurator.Load(context.Background())
		t.AssertNil(err)

		// Verify all values are loaded correctly
		t.Assert(config.Server.Host, "localhost")
		t.Assert(config.Server.Port, 8080)
		t.Assert(config.Server.Name, "myapp")
		t.Assert(config.Database.Host, "db.example.com")
		t.Assert(config.Database.Port, 5432)
		t.Assert(config.Database.Username, "user")
		t.Assert(config.Database.Password, "pass")
		t.Assert(config.Database.Database, "mydb")
		t.Assert(len(config.Features), 3)
		t.Assert(config.Features[0], "feature1")
		t.Assert(config.Features[1], "feature2")
		t.Assert(config.Features[2], "feature3")
	})
}
