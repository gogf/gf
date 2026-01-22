// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
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

// TestConfig2 is a test struct for configuration binding
type TestConfig2 struct {
	Name     string       `json:"name" yaml:"name"`
	Age      int          `json:"age" yaml:"age"`
	Enabled  bool         `json:"enabled" yaml:"enabled"`
	Features string       `json:"features" yaml:"features"`
	Server   ServerConfig `json:"server" yaml:"server"`
}

// TestConfig3 is a test struct for configuration binding
type TestConfig3 struct {
	Name     string       `json:"name" yaml:"name"`
	Age      int          `json:"age" yaml:"age"`
	Enabled  bool         `json:"enabled" yaml:"enabled"`
	Features []string     `json:"features" yaml:"features"`
	Server   ServerConfig `json:"server" yaml:"server"`
	Other    string       `json:"other" yaml:"other"`
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

func TestLoader_Load(t *testing.T) {
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

		// Create loader
		loader := gcfg.NewLoaderWithAdapter[TestConfig](cfg, "")

		// Load configuration
		err = loader.Load(context.Background())
		t.AssertNil(err)
		v := loader.Get()

		// Check loaded values
		t.Assert(v.Name, "test-app")
		t.Assert(v.Age, 25)
		t.Assert(v.Enabled, true)
		t.Assert(v.Server.Host, "localhost")
		t.Assert(v.Server.Port, 8080)
		t.Assert(len(v.Features), 3)
	})
}

func TestLoader_LoadWithDefaultValues(t *testing.T) {
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
		var targetConfig TestConfig3
		targetConfig.Other = "other"

		// Create loader
		loader := gcfg.NewLoaderWithAdapter(cfg, "", &targetConfig)
		loader.SetReuseTargetStruct(true)

		// Load configuration
		err = loader.Load(context.Background())
		t.AssertNil(err)
		v := loader.Get()

		// Check loaded values
		t.Assert(v.Name, "test-app")
		t.Assert(v.Age, 25)
		t.Assert(v.Enabled, true)
		t.Assert(v.Server.Host, "localhost")
		t.Assert(v.Server.Port, 8080)
		t.Assert(len(v.Features), 3)
		t.Assert(v.Other, "other")
	})
}

func TestLoader_LoadWithPropertyKey(t *testing.T) {
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

		// Create loader with specific property key
		loader := gcfg.NewLoaderWithAdapter[ServerConfig](cfg, "server")

		// Load configuration
		err = loader.Load(context.Background())
		t.AssertNil(err)
		v := loader.Get()

		// Check loaded values - only the app section should be loaded
		t.Assert(v.Host, "localhost")
		t.Assert(v.Port, 8080)
	})
}

func TestLoader_WatchAndOnChange(t *testing.T) {
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
		// Create a new config instance
		cfg, err := gcfg.NewAdapterContent(configContent)
		t.AssertNil(err)

		// Variable to track if callback was called
		callbackCalled := gtype.NewBool(false)

		// Create loader
		loader := gcfg.NewLoaderWithAdapter[TestConfig](cfg, "")

		// Set change callback
		loader.OnChange(func(updated TestConfig) error {
			callbackCalled.Set(true)
			return nil
		})

		// Load configuration
		err = loader.Load(context.Background())
		t.AssertNil(err)
		err = loader.Watch(context.Background(), "test-watcher")
		t.AssertNil(err)
		v := loader.Get()
		t.Assert(v.Name, "test-app")
		t.Assert(v.Age, 25)
		err = cfg.SetContent(configContent2)
		t.AssertNil(err)
		time.Sleep(2 * time.Second)
		v2 := loader.Get()
		t.Assert(v2.Name, "test-app-2")
		t.Assert(v2.Age, 200)
		t.Assert(callbackCalled.Val(), true)
	})
}

func TestLoader_SetConverter(t *testing.T) {
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

		// Create loader
		loader := gcfg.NewLoaderWithAdapter[TestConfig2](cfg, "features")

		// Set custom converter
		loader.SetConverter(func(data any, target *TestConfig2) error {
			s := gconv.Strings(data)
			target.Features = strings.Join(s, ",")
			return nil
		})

		// Load configuration
		err = loader.Load(context.Background())
		t.AssertNil(err)
		v := loader.Get()

		// Check converted values
		t.Assert(v.Features, "feature,feature,feature")
	})
}

func TestLoader_SetWatchErrorHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a new config instance with content that will cause converter error
		cfg, err := gcfg.NewAdapterContent(configContent)
		t.AssertNil(err)

		// Create loader
		loader := gcfg.NewLoaderWithAdapter[TestConfig](cfg, "")

		// Set error handler for watch operations
		errorHandled := gtype.NewBool(false)
		loader.SetWatchErrorHandler(func(ctx context.Context, err error) {
			errorHandled.Set(true)
		})

		// Set a converter that will fail
		loader.SetConverter(func(data any, target *TestConfig) error {
			return errors.New("converter error")
		})

		// Load initially - this should return error without calling error handler
		err = loader.Load(context.Background())
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "converter error")
		// Error handler should NOT be called during direct Load
		t.Assert(errorHandled.Val(), false)

		// Start watching - now errors during Load should trigger the error handler
		err = loader.Watch(context.Background(), "test-error-handler")
		t.AssertNil(err)
		// Reset
		errorHandled.Set(false)
		// Trigger a config change - this will call Load internally and should trigger error handler
		err = cfg.SetContent(configContent)
		t.AssertNil(err)

		// Wait for watcher to process the change
		time.Sleep(1 * time.Second)

		// Error handler should be called during Watch's Load
		t.Assert(errorHandled.Val(), true)
	})
}

func TestLoader_IsWatchingAndStopWatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a new config instance
		cfg, err := gcfg.NewAdapterContent(configContent)
		t.AssertNil(err)

		// Create loader
		loader := gcfg.NewLoaderWithAdapter[TestConfig](cfg, "")

		// Initially, should not be watching
		t.Assert(loader.IsWatching(), false)

		// Load configuration
		err = loader.Load(context.Background())
		t.AssertNil(err)

		// Start watching
		err = loader.Watch(context.Background(), "test-stopwatch-watcher")
		t.AssertNil(err)

		// Now should be watching
		t.Assert(loader.IsWatching(), true)

		// Stop watching
		stopped, err := loader.StopWatch(context.Background())
		t.AssertNil(err)
		t.Assert(stopped, true)

		// Should not be watching anymore
		t.Assert(loader.IsWatching(), false)
	})
}

func TestLoader_StopWatchWithoutWatcher(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create a new config instance
		cfg, err := gcfg.NewAdapterContent(configContent)
		t.AssertNil(err)

		// Create loader without starting to watch
		loader := gcfg.NewLoaderWithAdapter[TestConfig](cfg, "")

		// Initially, should not be watching
		t.Assert(loader.IsWatching(), false)

		// Try to stop watching when not watching
		stopped, err := loader.StopWatch(context.Background())
		t.AssertNE(err, nil)
		t.Assert(stopped, false)
		t.Assert(err.Error(), "No watcher name specified")
	})
}
