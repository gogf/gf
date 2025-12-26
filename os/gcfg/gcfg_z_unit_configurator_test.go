// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg_test

import (
	"context"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// TestConfig is a test struct for configuration binding
type TestConfig struct {
	Name     string       `json:"name" yaml:"name"`
	Age      int          `json:"age" yaml:"age"`
	Enabled  bool         `json:"enabled" yaml:"enabled"`
	Features []string     `json:"features" yaml:"features"`
	Server   ServerConfig `json:"server" yaml:"server"`
}

type TestConfig2 struct {
	Name     string       `json:"name" yaml:"name"`
	Age      int          `json:"age" yaml:"age"`
	Enabled  bool         `json:"enabled" yaml:"enabled"`
	Features string       `json:"features" yaml:"features"`
	Server   ServerConfig `json:"server" yaml:"server"`
}

type ServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

var configContent = `
name: "test-app"
age: 25
enabled: true
features: ["feature1", "feature2", "feature3"]
server:
  host: "localhost"
  port: 8080
`

func TestConfigurator_Load(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			configFile = "./" + guid.S() + ".yaml"
			err        = gfile.PutContents(configFile, configContent)
		)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create a new config instance
		cfg, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)

		// Create target struct
		var targetConfig TestConfig

		// Create configurator
		configurator := gcfg.NewConfiguratorWithAdapter[TestConfig](cfg, "", &targetConfig)

		// Load configuration
		err = configurator.Load(context.Background())
		t.AssertNil(err)
		v := configurator.Get()

		// Check loaded values
		t.Assert(v.Name, "test-app")
		t.Assert(v.Age, 25)
		t.Assert(v.Enabled, true)
		t.Assert(v.Server.Host, "localhost")
		t.Assert(v.Server.Port, 8080)
		t.Assert(len(v.Features), 3)
	})
}

func TestConfigurator_LoadWithPropertyKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			configFile = "./" + guid.S() + ".yaml"
			err        = gfile.PutContents(configFile, configContent)
		)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create a new config instance
		cfg, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)

		// Create target struct
		var targetConfig ServerConfig

		// Create configurator with specific property key
		configurator := gcfg.NewConfiguratorWithAdapter[ServerConfig](cfg, "server", &targetConfig)

		// Load configuration
		err = configurator.Load(context.Background())
		t.AssertNil(err)
		v := configurator.Get()

		// Check loaded values - only the app section should be loaded
		t.Assert(v.Host, "localhost")
		t.Assert(v.Port, 8080)
	})
}

func TestConfigurator_WatchAndOnChange(t *testing.T) {
	var configContent2 = `
name: test-app-2
age: 200
enabled: true
features: ["feature1", "feature2", "feature3"]
server:
  host: localhost
  port: 8080
`

	gtest.C(t, func(t *gtest.T) {
		var (
			configFile = "./" + guid.S() + ".yaml"
			err        = gfile.PutContents(configFile, configContent)
		)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create a new config instance
		cfg, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)

		// Create target struct
		var targetConfig TestConfig

		// Variable to track if callback was called
		callbackCalled := gtype.NewBool(false)

		// Create configurator
		configurator := gcfg.NewConfiguratorWithAdapter[TestConfig](cfg, "", &targetConfig)

		// Set change callback
		configurator.OnChange(func(updated TestConfig) error {
			callbackCalled.Set(true)
			return nil
		})

		// Load configuration
		err = configurator.Load(context.Background())
		t.AssertNil(err)
		err = configurator.Watch(context.Background(), "test-watcher")
		t.AssertNil(err)
		v := configurator.Get()
		t.Assert(v.Name, "test-app")
		t.Assert(v.Age, 25)
		err = gfile.PutContents(configFile, configContent2)
		t.AssertNil(err)
		time.Sleep(1 * time.Second)
		v2 := configurator.Get()
		t.Assert(v2.Name, "test-app-2")
		t.Assert(v2.Age, 200)
		t.Assert(callbackCalled.Val(), true)
	})
}

func TestConfigurator_SetConverter(t *testing.T) {
	var configContent2 = `
name: test-app-2
age: 200
enabled: true
features: ["feature", "feature", "feature"]
server:
  host: localhost
  port: 8080
`
	gtest.C(t, func(t *gtest.T) {
		var (
			configFile = "./" + guid.S() + ".yaml"
			err        = gfile.PutContents(configFile, configContent2)
		)
		t.AssertNil(err)
		defer gfile.RemoveFile(configFile)

		// Create a new config instance
		cfg, err := gcfg.NewAdapterFile(configFile)
		t.AssertNil(err)

		// Create target struct
		var targetConfig TestConfig2

		// Create configurator
		configurator := gcfg.NewConfiguratorWithAdapter[TestConfig2](cfg, "features", &targetConfig)

		// Set custom converter
		configurator.SetConverter(func(data any, target *TestConfig2) error {
			s := gconv.Strings(data)
			target.Features = strings.Join(s, ",")
			return nil
		})

		// Load configuration
		err = configurator.Load(context.Background())
		t.AssertNil(err)
		v := configurator.Get()

		// Check converted values
		t.Assert(v.Features, "feature,feature,feature")
	})
}

func TestConfigurator_SetLoadErrorHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a new config instance with invalid adapter that will cause error
		cfg, err := gcfg.New()
		t.AssertNil(err)

		// Create target struct
		var targetConfig TestConfig

		// Create configurator
		configurator := gcfg.NewConfigurator(cfg, "non-existent-key", &targetConfig)

		// Set error handler
		errorHandled := gtype.NewBool(false)
		configurator.SetLoadErrorHandler(func(ctx context.Context, err error) {
			errorHandled.Set(true)
		})

		// Try to load with invalid key
		err = configurator.Load(context.Background())
		// The error should be handled by our error handler
		t.Assert(errorHandled, true)
	})
}
