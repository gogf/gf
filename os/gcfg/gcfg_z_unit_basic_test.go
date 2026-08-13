// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Basic1(t *testing.T) {
	config := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gtest.C(t, func(t *gtest.T) {
		var (
			path = gcfg.DefaultConfigFileName
			err  = gfile.PutContents(path, config)
		)
		t.AssertNil(err)
		defer gfile.Remove(path)

		c, err := gcfg.New()
		t.AssertNil(err)
		t.Assert(c.MustGet(ctx, "v1"), 1)
		filepath, _ := c.GetAdapter().(*gcfg.AdapterFile).GetFilePath()
		t.AssertEQ(filepath, gfile.Pwd()+gfile.Separator+path)
	})
}

func Test_Basic2(t *testing.T) {
	config := `log-path = "logs"`
	gtest.C(t, func(t *gtest.T) {
		var (
			path = gcfg.DefaultConfigFileName
			err  = gfile.PutContents(path, config)
		)
		t.AssertNil(err)
		defer func() {
			_ = gfile.Remove(path)
		}()

		c, err := gcfg.New()
		t.AssertNil(err)
		t.Assert(c.MustGet(ctx, "log-path"), "logs")
	})
}

func Test_Content(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()
		t.Assert(c.MustGet(ctx, "v1"), 1)
	})
}

func Test_SetFileName(t *testing.T) {
	config := `
{
	"array": [
		1,
		2,
		3
	],
	"redis": {
		"cache": "127.0.0.1:6379,1",
		"disk": "127.0.0.1:6379,0"
	},
	"v1": 1,
	"v2": "true",
	"v3": "off",
	"v4": "1.234"
}
`
	gtest.C(t, func(t *gtest.T) {
		path := "config.json"
		err := gfile.PutContents(path, config)
		t.AssertNil(err)
		defer func() {
			_ = gfile.Remove(path)
		}()

		config, err := gcfg.New()
		t.AssertNil(err)
		c := config.GetAdapter().(*gcfg.AdapterFile)
		c.SetFileName(path)
		t.Assert(c.MustGet(ctx, "v1"), 1)
		t.AssertEQ(c.MustGet(ctx, "v1").Int(), 1)
		t.AssertEQ(c.MustGet(ctx, "v1").Int8(), int8(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Int16(), int16(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Int32(), int32(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Int64(), int64(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint(), uint(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint8(), uint8(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint16(), uint16(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint32(), uint32(1))
		t.AssertEQ(c.MustGet(ctx, "v1").Uint64(), uint64(1))

		t.AssertEQ(c.MustGet(ctx, "v1").String(), "1")
		t.AssertEQ(c.MustGet(ctx, "v1").Bool(), true)
		t.AssertEQ(c.MustGet(ctx, "v2").String(), "true")
		t.AssertEQ(c.MustGet(ctx, "v2").Bool(), true)

		t.AssertEQ(c.MustGet(ctx, "v1").String(), "1")
		t.AssertEQ(c.MustGet(ctx, "v4").Float32(), float32(1.234))
		t.AssertEQ(c.MustGet(ctx, "v4").Float64(), float64(1.234))
		t.AssertEQ(c.MustGet(ctx, "v2").String(), "true")
		t.AssertEQ(c.MustGet(ctx, "v2").Bool(), true)
		t.AssertEQ(c.MustGet(ctx, "v3").Bool(), false)

		t.AssertEQ(c.MustGet(ctx, "array").Ints(), []int{1, 2, 3})
		t.AssertEQ(c.MustGet(ctx, "array").Strings(), []string{"1", "2", "3"})
		t.AssertEQ(c.MustGet(ctx, "array").Interfaces(), []any{1, 2, 3})
		t.AssertEQ(c.MustGet(ctx, "redis").Map(), map[string]any{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		filepath, _ := c.GetFilePath()
		t.AssertEQ(filepath, gfile.Pwd()+gfile.Separator+path)
	})
}

func TestCfg_Get_WrongConfigFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var err error
		configPath := gfile.Temp(gtime.TimestampNanoStr())
		err = gfile.Mkdir(configPath)
		t.AssertNil(err)
		defer gfile.Remove(configPath)

		defer gfile.Chdir(gfile.Pwd())
		err = gfile.Chdir(configPath)
		t.AssertNil(err)

		err = gfile.PutContents(
			gfile.Join(configPath, "config.yml"),
			"wrong config",
		)
		t.AssertNil(err)
		adapterFile, err := gcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		c := gcfg.NewWithAdapter(adapterFile)
		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
		adapterFile.Clear()
	})
}

func Test_GetWithEnv(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()
		t.Assert(c.MustGet(ctx, "v1"), 1)
		t.Assert(c.MustGetWithEnv(ctx, `redis.user`), nil)
		t.Assert(genv.Set("REDIS_USER", `1`), nil)
		defer genv.Remove(`REDIS_USER`)
		t.Assert(c.MustGetWithEnv(ctx, `redis.user`), `1`)
	})
}

func Test_GetWithCmd(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gtest.C(t, func(t *gtest.T) {

		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()
		t.Assert(c.MustGet(ctx, "v1"), 1)
		t.Assert(c.MustGetWithCmd(ctx, `redis.user`), nil)

		gcmd.Init([]string{"gf", "--redis.user=2"}...)
		t.Assert(c.MustGetWithCmd(ctx, `redis.user`), `2`)
	})
}

func Test_GetEffective(t *testing.T) {
	content := `
v1    = 1
v2    = "true"
[server]
    port = 8080
    host = "localhost"
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()

		// Test 1: Get from config file when no cmd/env set
		t.Assert(c.MustGetEffective(ctx, "server.port"), 8080)
		t.Assert(c.MustGetEffective(ctx, "server.host"), "localhost")

		// Test 2: Environment variable overrides config file
		t.Assert(genv.Set("SERVER_PORT", "9090"), nil)
		defer genv.Remove("SERVER_PORT")
		t.Assert(c.MustGetEffective(ctx, "server.port"), "9090")

		// Test 3: Command line overrides environment variable
		gcmd.Init([]string{"gf", "--server.port=7070"}...)
		t.Assert(c.MustGetEffective(ctx, "server.port"), "7070")

		// Test 4: Default value when nothing is set
		t.Assert(c.MustGetEffective(ctx, "server.timeout", 30), 30)

		// Test 5: Empty string from command line should override
		gcmd.Init([]string{"gf", "--server.name="}...)
		t.Assert(genv.Set("SERVER_NAME", "from-env"), nil)
		defer genv.Remove("SERVER_NAME")
		t.Assert(c.MustGetEffective(ctx, "server.name"), "")

		// Test 6: Key not in config, only in env
		t.Assert(genv.Set("APP_DEBUG", "true"), nil)
		defer genv.Remove("APP_DEBUG")
		t.Assert(c.MustGetEffective(ctx, "app.debug"), "true")
	})
}

func Test_GetWithEnv_NormalizePattern(t *testing.T) {
	// Test configuration where keys with underscores should be normalized to dots
	content := `
[server]
    address = "127.0.0.1:8000"
    port = 8080
[redis]
    user = "default"
`
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()

		// Test 1: GetWithEnv with uppercase underscore pattern should find config with dot notation.
		// SERVER_ADDRESS should be normalized to server.address and find the config value.
		t.Assert(c.MustGetWithEnv(ctx, "SERVER_ADDRESS").String(), "127.0.0.1:8000")

		// Test 2: GetWithEnv with dot notation pattern should still work (regression test).
		t.Assert(c.MustGetWithEnv(ctx, "server.address").String(), "127.0.0.1:8000")

		// Test 3: GetWithEnv with uppercase underscore pattern for existing config.
		t.Assert(c.MustGetWithEnv(ctx, "SERVER_PORT").Int(), 8080)

		// Test 4: GetWithEnv with dot notation pattern for existing config.
		t.Assert(c.MustGetWithEnv(ctx, "server.port").Int(), 8080)

		// Test 5: When config exists, env variable should not override (config takes priority).
		t.Assert(genv.Set("SERVER_ADDRESS", "from-env"), nil)
		defer genv.Remove("SERVER_ADDRESS")
		// Config value should take precedence over env variable.
		t.Assert(c.MustGetWithEnv(ctx, "SERVER_ADDRESS").String(), "127.0.0.1:8000")
		t.Assert(c.MustGetWithEnv(ctx, "server.address").String(), "127.0.0.1:8000")

		// Test 6: When config doesn't exist, should fallback to env variable.
		t.Assert(genv.Set("APP_NAME", "myapp"), nil)
		defer genv.Remove("APP_NAME")
		// APP_NAME should be normalized to app.name, not found in config, fallback to env.
		t.Assert(c.MustGetWithEnv(ctx, "APP_NAME").String(), "myapp")

		// Test 7: Mixed case with underscores fallback to env.
		t.Assert(genv.Set("MY_CUSTOM_KEY", "custom-value"), nil)
		defer genv.Remove("MY_CUSTOM_KEY")
		// MY_CUSTOM_KEY should be normalized to my.custom.key, not found in config, fallback to env.
		t.Assert(c.MustGetWithEnv(ctx, "MY_CUSTOM_KEY").String(), "custom-value")
	})
}

// Test_FormatCmdKey_BehaviorExplanation documents the boundary behavior of FormatCmdKey conversion.
// FormatCmdKey converts to lowercase and replaces underscores with dots.
//
// Conversion rules:
//   - "my_custom_key" → "my.custom.key"
//   - "SERVER_ADDRESS" → "server.address"
//   - "a.b_c.d" → "a.b.c.d"
//
// Note: The pattern "_" converts to "." which is a special pattern in config
// that retrieves all configuration data.
func Test_FormatCmdKey_BehaviorExplanation(t *testing.T) {
	// Test configuration with various key formats that correspond to FormatCmdKey conversion results
	content := `
[my.custom]
    key = "value1"
[a.b.c]
    d = "value2"
[server]
    address = "127.0.0.1:8000"
`
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()

		// Test 1: Simple underscore conversion: MY_CUSTOM_KEY → my.custom.key
		// FormatCmdKey("MY_CUSTOM_KEY") = "my.custom.key"
		// Config has [my] custom_key = "value1", which matches "my.custom.key" after normalization
		t.Assert(c.MustGetWithEnv(ctx, "MY_CUSTOM_KEY").String(), "value1")

		// Test 2: Mixed dot and underscore: A_B_C_D → a.b.c.d (FormatCmdKey converts underscores to dots)
		// FormatCmdKey("A_B_C_D") = "a.b.c.d"
		// Config has [a] b_c_d = "value2", which matches "a.b.c.d" after normalization
		t.Assert(c.MustGetWithEnv(ctx, "A_B_C_D").String(), "value2")

		// Test 3: Pure dot notation: server.address should still work
		t.Assert(c.MustGetWithEnv(ctx, "server.address").String(), "127.0.0.1:8000")

		// Test 4: Uppercase with dots: SERVER.ADDRESS → FormatCmdKey("SERVER.ADDRESS") = "server.address"
		// Config has server.address, so it should match
		t.Assert(c.MustGetWithEnv(ctx, "SERVER.ADDRESS").String(), "127.0.0.1:8000")

		// Test 5: All underscores: ALL_UNDERSCORES → all.underscores (no matching config, fallback to env)
		t.Assert(genv.Set("ALL_UNDERSCORES", "underscore-value"), nil)
		defer genv.Remove("ALL_UNDERSCORES")
		t.Assert(c.MustGetWithEnv(ctx, "ALL_UNDERSCORES").String(), "underscore-value")

		// Test 6: Single character key (no matching config, fallback to env)
		t.Assert(genv.Set("K", "single"), nil)
		defer genv.Remove("K")
		t.Assert(c.MustGetWithEnv(ctx, "K").String(), "single")
	})
}

func Test_GetWithCmd_NormalizePattern(t *testing.T) {
	content := `
[server]
    address = "127.0.0.1:8000"
    port = 8080
`
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()

		// Reset command line state
		gcmd.Init([]string{"gf"}...)

		// Test 1: GetWithCmd with dot notation pattern should find config.
		t.Assert(c.MustGetWithCmd(ctx, "server.address").String(), "127.0.0.1:8000")

		// Test 2: GetWithCmd with underscore pattern should be normalized and find config.
		t.Assert(c.MustGetWithCmd(ctx, "SERVER_ADDRESS").String(), "127.0.0.1:8000")

		// Test 3: GetWithCmd with dot notation for port.
		t.Assert(c.MustGetWithCmd(ctx, "server.port").Int(), 8080)

		// Test 4: GetWithCmd with underscore pattern for port.
		t.Assert(c.MustGetWithCmd(ctx, "SERVER_PORT").Int(), 8080)

		// Test 5: Command line should be used as fallback when config value not found.
		gcmd.Init([]string{"gf", "--app.name=myapp"}...)
		t.Assert(c.MustGetWithCmd(ctx, "APP_NAME").String(), "myapp")
	})
}

func Test_GetEffective_NormalizePattern(t *testing.T) {
	content := `
[server]
    address = "127.0.0.1:8000"
    port = 8080
`
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.New()
		t.AssertNil(err)
		c.GetAdapter().(*gcfg.AdapterFile).SetContent(content)
		defer c.GetAdapter().(*gcfg.AdapterFile).ClearContent()

		// Reset command line state
		gcmd.Init([]string{"gf"}...)

		// Test 1: GetEffective with underscore pattern should find config via normalized key.
		t.Assert(c.MustGetEffective(ctx, "SERVER_ADDRESS").String(), "127.0.0.1:8000")

		// Test 2: GetEffective with dot notation should still work (regression).
		t.Assert(c.MustGetEffective(ctx, "server.address").String(), "127.0.0.1:8000")

		// Test 3: Environment variable should override config in GetEffective.
		t.Assert(genv.Set("SERVER_ADDRESS", "from-env"), nil)
		defer genv.Remove("SERVER_ADDRESS")
		t.Assert(c.MustGetEffective(ctx, "SERVER_ADDRESS").String(), "from-env")

		// Test 4: Command line should override both env and config.
		gcmd.Init([]string{"gf", "--server.address=from-cmd"}...)
		t.Assert(c.MustGetEffective(ctx, "SERVER_ADDRESS").String(), "from-cmd")
	})
}
