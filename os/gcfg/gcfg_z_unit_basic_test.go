// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
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
		path := gcfg.DefaultConfigFile
		err := gfile.PutContents(path, config)
		t.Assert(err, nil)
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
		path := gcfg.DefaultConfigFile
		err := gfile.PutContents(path, config)
		t.Assert(err, nil)
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
		t.Assert(err, nil)
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
		t.AssertEQ(c.MustGet(ctx, "array").Interfaces(), []interface{}{1, 2, 3})
		t.AssertEQ(c.MustGet(ctx, "redis").Map(), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		filepath, _ := c.GetFilePath()
		t.AssertEQ(filepath, gfile.Pwd()+gfile.Separator+path)
	})
}

func TestCfg_Set(t *testing.T) {
	config := `log-path = "logs"`
	gtest.C(t, func(t *gtest.T) {
		path := gcfg.DefaultConfigFile
		err := gfile.PutContents(path, config)
		t.Assert(err, nil)
		defer gfile.Remove(path)

		adapterFile, err := gcfg.NewAdapterFile()
		t.AssertNil(err)
		t.Assert(adapterFile.MustGet(ctx, "log-path"), "logs")

		c := gcfg.NewWithAdapter(adapterFile)
		c.Set(ctx, "log-path", "custom-logs")
		t.Assert(err, nil)
		t.Assert(c.MustGet(ctx, "log-path"), "custom-logs")
	})
}

func TestCfg_Get_WrongConfigFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var err error
		configPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(configPath)
		t.Assert(err, nil)
		defer gfile.Remove(configPath)

		defer gfile.Chdir(gfile.Pwd())
		err = gfile.Chdir(configPath)
		t.Assert(err, nil)

		err = gfile.PutContents(
			gfile.Join(configPath, "config.yml"),
			"wrong config",
		)
		t.Assert(err, nil)
		adapterFile, err := gcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		c := gcfg.NewWithAdapter(adapterFile)
		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
		adapterFile.Clear()
	})
}
