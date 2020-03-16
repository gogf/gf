// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg_test

import (
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/gtime"
	"os"
	"testing"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func init() {
	os.Setenv("GF_GCFG_ERRORPRINT", "false")
}

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
	gtest.Case(t, func() {
		path := gcfg.DEFAULT_CONFIG_FILE
		err := gfile.PutContents(path, config)
		gtest.Assert(err, nil)
		defer gfile.Remove(path)

		c := gcfg.New()
		gtest.Assert(c.Get("v1"), 1)
		gtest.AssertEQ(c.GetInt("v1"), 1)
		gtest.AssertEQ(c.GetInt8("v1"), int8(1))
		gtest.AssertEQ(c.GetInt16("v1"), int16(1))
		gtest.AssertEQ(c.GetInt32("v1"), int32(1))
		gtest.AssertEQ(c.GetInt64("v1"), int64(1))
		gtest.AssertEQ(c.GetUint("v1"), uint(1))
		gtest.AssertEQ(c.GetUint8("v1"), uint8(1))
		gtest.AssertEQ(c.GetUint16("v1"), uint16(1))
		gtest.AssertEQ(c.GetUint32("v1"), uint32(1))
		gtest.AssertEQ(c.GetUint64("v1"), uint64(1))

		gtest.AssertEQ(c.GetVar("v1").String(), "1")
		gtest.AssertEQ(c.GetVar("v1").Bool(), true)
		gtest.AssertEQ(c.GetVar("v2").String(), "true")
		gtest.AssertEQ(c.GetVar("v2").Bool(), true)

		gtest.AssertEQ(c.GetString("v1"), "1")
		gtest.AssertEQ(c.GetFloat32("v4"), float32(1.23))
		gtest.AssertEQ(c.GetFloat64("v4"), float64(1.23))
		gtest.AssertEQ(c.GetString("v2"), "true")
		gtest.AssertEQ(c.GetBool("v2"), true)
		gtest.AssertEQ(c.GetBool("v3"), false)

		gtest.AssertEQ(c.Contains("v1"), true)
		gtest.AssertEQ(c.Contains("v2"), true)
		gtest.AssertEQ(c.Contains("v3"), true)
		gtest.AssertEQ(c.Contains("v4"), true)
		gtest.AssertEQ(c.Contains("v5"), false)

		gtest.AssertEQ(c.GetInts("array"), []int{1, 2, 3})
		gtest.AssertEQ(c.GetStrings("array"), []string{"1", "2", "3"})
		gtest.AssertEQ(c.GetArray("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetInterfaces("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetMap("redis"), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		gtest.AssertEQ(c.FilePath(), gfile.Pwd()+gfile.Separator+path)

	})
}

func Test_Basic2(t *testing.T) {
	config := `log-path = "logs"`
	gtest.Case(t, func() {
		path := gcfg.DEFAULT_CONFIG_FILE
		err := gfile.PutContents(path, config)
		gtest.Assert(err, nil)
		defer func() {
			_ = gfile.Remove(path)
		}()

		c := gcfg.New()
		gtest.Assert(c.Get("log-path"), "logs")
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
	gcfg.SetContent(content)
	defer gcfg.ClearContent()

	gtest.Case(t, func() {
		c := gcfg.New()
		gtest.Assert(c.Get("v1"), 1)
		gtest.AssertEQ(c.GetInt("v1"), 1)
		gtest.AssertEQ(c.GetInt8("v1"), int8(1))
		gtest.AssertEQ(c.GetInt16("v1"), int16(1))
		gtest.AssertEQ(c.GetInt32("v1"), int32(1))
		gtest.AssertEQ(c.GetInt64("v1"), int64(1))
		gtest.AssertEQ(c.GetUint("v1"), uint(1))
		gtest.AssertEQ(c.GetUint8("v1"), uint8(1))
		gtest.AssertEQ(c.GetUint16("v1"), uint16(1))
		gtest.AssertEQ(c.GetUint32("v1"), uint32(1))
		gtest.AssertEQ(c.GetUint64("v1"), uint64(1))

		gtest.AssertEQ(c.GetVar("v1").String(), "1")
		gtest.AssertEQ(c.GetVar("v1").Bool(), true)
		gtest.AssertEQ(c.GetVar("v2").String(), "true")
		gtest.AssertEQ(c.GetVar("v2").Bool(), true)

		gtest.AssertEQ(c.GetString("v1"), "1")
		gtest.AssertEQ(c.GetFloat32("v4"), float32(1.23))
		gtest.AssertEQ(c.GetFloat64("v4"), float64(1.23))
		gtest.AssertEQ(c.GetString("v2"), "true")
		gtest.AssertEQ(c.GetBool("v2"), true)
		gtest.AssertEQ(c.GetBool("v3"), false)

		gtest.AssertEQ(c.Contains("v1"), true)
		gtest.AssertEQ(c.Contains("v2"), true)
		gtest.AssertEQ(c.Contains("v3"), true)
		gtest.AssertEQ(c.Contains("v4"), true)
		gtest.AssertEQ(c.Contains("v5"), false)

		gtest.AssertEQ(c.GetInts("array"), []int{1, 2, 3})
		gtest.AssertEQ(c.GetStrings("array"), []string{"1", "2", "3"})
		gtest.AssertEQ(c.GetArray("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetInterfaces("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetMap("redis"), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
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
	gtest.Case(t, func() {
		path := "config.json"
		err := gfile.PutContents(path, config)
		gtest.Assert(err, nil)
		defer func() {
			_ = gfile.Remove(path)
		}()

		c := gcfg.New()
		c.SetFileName(path)
		gtest.Assert(c.Get("v1"), 1)
		gtest.AssertEQ(c.GetInt("v1"), 1)
		gtest.AssertEQ(c.GetInt8("v1"), int8(1))
		gtest.AssertEQ(c.GetInt16("v1"), int16(1))
		gtest.AssertEQ(c.GetInt32("v1"), int32(1))
		gtest.AssertEQ(c.GetInt64("v1"), int64(1))
		gtest.AssertEQ(c.GetUint("v1"), uint(1))
		gtest.AssertEQ(c.GetUint8("v1"), uint8(1))
		gtest.AssertEQ(c.GetUint16("v1"), uint16(1))
		gtest.AssertEQ(c.GetUint32("v1"), uint32(1))
		gtest.AssertEQ(c.GetUint64("v1"), uint64(1))

		gtest.AssertEQ(c.GetVar("v1").String(), "1")
		gtest.AssertEQ(c.GetVar("v1").Bool(), true)
		gtest.AssertEQ(c.GetVar("v2").String(), "true")
		gtest.AssertEQ(c.GetVar("v2").Bool(), true)

		gtest.AssertEQ(c.GetString("v1"), "1")
		gtest.AssertEQ(c.GetFloat32("v4"), float32(1.234))
		gtest.AssertEQ(c.GetFloat64("v4"), float64(1.234))
		gtest.AssertEQ(c.GetString("v2"), "true")
		gtest.AssertEQ(c.GetBool("v2"), true)
		gtest.AssertEQ(c.GetBool("v3"), false)

		gtest.AssertEQ(c.Contains("v1"), true)
		gtest.AssertEQ(c.Contains("v2"), true)
		gtest.AssertEQ(c.Contains("v3"), true)
		gtest.AssertEQ(c.Contains("v4"), true)
		gtest.AssertEQ(c.Contains("v5"), false)

		gtest.AssertEQ(c.GetInts("array"), []int{1, 2, 3})
		gtest.AssertEQ(c.GetStrings("array"), []string{"1", "2", "3"})
		gtest.AssertEQ(c.GetArray("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetInterfaces("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetMap("redis"), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		gtest.AssertEQ(c.FilePath(), gfile.Pwd()+gfile.Separator+path)

	})
}

func Test_Instance(t *testing.T) {
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
	gtest.Case(t, func() {
		path := gcfg.DEFAULT_CONFIG_FILE
		err := gfile.PutContents(path, config)
		gtest.Assert(err, nil)
		defer func() {
			gtest.Assert(gfile.Remove(path), nil)
		}()

		c := gcfg.Instance()
		gtest.Assert(c.Get("v1"), 1)
		gtest.AssertEQ(c.GetInt("v1"), 1)
		gtest.AssertEQ(c.GetInt8("v1"), int8(1))
		gtest.AssertEQ(c.GetInt16("v1"), int16(1))
		gtest.AssertEQ(c.GetInt32("v1"), int32(1))
		gtest.AssertEQ(c.GetInt64("v1"), int64(1))
		gtest.AssertEQ(c.GetUint("v1"), uint(1))
		gtest.AssertEQ(c.GetUint8("v1"), uint8(1))
		gtest.AssertEQ(c.GetUint16("v1"), uint16(1))
		gtest.AssertEQ(c.GetUint32("v1"), uint32(1))
		gtest.AssertEQ(c.GetUint64("v1"), uint64(1))

		gtest.AssertEQ(c.GetVar("v1").String(), "1")
		gtest.AssertEQ(c.GetVar("v1").Bool(), true)
		gtest.AssertEQ(c.GetVar("v2").String(), "true")
		gtest.AssertEQ(c.GetVar("v2").Bool(), true)

		gtest.AssertEQ(c.GetString("v1"), "1")
		gtest.AssertEQ(c.GetFloat32("v4"), float32(1.234))
		gtest.AssertEQ(c.GetFloat64("v4"), float64(1.234))
		gtest.AssertEQ(c.GetString("v2"), "true")
		gtest.AssertEQ(c.GetBool("v2"), true)
		gtest.AssertEQ(c.GetBool("v3"), false)

		gtest.AssertEQ(c.Contains("v1"), true)
		gtest.AssertEQ(c.Contains("v2"), true)
		gtest.AssertEQ(c.Contains("v3"), true)
		gtest.AssertEQ(c.Contains("v4"), true)
		gtest.AssertEQ(c.Contains("v5"), false)

		gtest.AssertEQ(c.GetInts("array"), []int{1, 2, 3})
		gtest.AssertEQ(c.GetStrings("array"), []string{"1", "2", "3"})
		gtest.AssertEQ(c.GetArray("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetInterfaces("array"), []interface{}{1, 2, 3})
		gtest.AssertEQ(c.GetMap("redis"), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		gtest.AssertEQ(c.FilePath(), gfile.Pwd()+gfile.Separator+path)

	})
}

func TestCfg_New(t *testing.T) {
	gtest.Case(t, func() {
		os.Setenv("GF_GCFG_PATH", "config")
		c := gcfg.New("config.yml")
		gtest.Assert(c.Get("name"), nil)
		gtest.Assert(c.GetFileName(), "config.yml")

		configPath := gfile.Pwd() + gfile.Separator + "config"
		_ = gfile.Mkdir(configPath)
		defer gfile.Remove(configPath)

		c = gcfg.New("config.yml")
		gtest.Assert(c.Get("name"), nil)

		_ = os.Unsetenv("GF_GCFG_PATH")
		c = gcfg.New("config.yml")
		gtest.Assert(c.Get("name"), nil)
	})
}

func TestCfg_SetPath(t *testing.T) {
	gtest.Case(t, func() {
		c := gcfg.New("config.yml")
		err := c.SetPath("tmp")
		gtest.AssertNE(err, nil)
		err = c.SetPath("gcfg.go")
		gtest.AssertNE(err, nil)
		gtest.Assert(c.Get("name"), nil)
	})
}

func TestCfg_SetViolenceCheck(t *testing.T) {
	gtest.Case(t, func() {
		c := gcfg.New("config.yml")
		c.SetViolenceCheck(true)
		gtest.Assert(c.Get("name"), nil)
	})
}

func TestCfg_AddPath(t *testing.T) {
	gtest.Case(t, func() {
		c := gcfg.New("config.yml")
		err := c.AddPath("tmp")
		gtest.AssertNE(err, nil)
		err = c.AddPath("gcfg.go")
		gtest.AssertNE(err, nil)
		gtest.Assert(c.Get("name"), nil)
	})
}

func TestCfg_FilePath(t *testing.T) {
	gtest.Case(t, func() {
		c := gcfg.New("config.yml")
		path := c.FilePath("tmp")
		gtest.Assert(path, "")
		path = c.FilePath("tmp")
		gtest.Assert(path, "")
	})
}

func TestCfg_Get(t *testing.T) {
	gtest.Case(t, func() {
		var err error
		configPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err = gfile.Mkdir(configPath)
		gtest.Assert(err, nil)
		defer gfile.Remove(configPath)

		defer gfile.Chdir(gfile.Pwd())
		err = gfile.Chdir(configPath)
		gtest.Assert(err, nil)

		err = gfile.PutContents(
			gfile.Join(configPath, "config.yml"),
			"wrong config",
		)
		gtest.Assert(err, nil)
		c := gcfg.New("config.yml")
		gtest.Assert(c.Get("name"), nil)
		gtest.Assert(c.GetVar("name").Val(), nil)
		gtest.Assert(c.Contains("name"), false)
		gtest.Assert(c.GetMap("name"), nil)
		gtest.Assert(c.GetArray("name"), nil)
		gtest.Assert(c.GetString("name"), "")
		gtest.Assert(c.GetStrings("name"), nil)
		gtest.Assert(c.GetInterfaces("name"), nil)
		gtest.Assert(c.GetBool("name"), false)
		gtest.Assert(c.GetFloat32("name"), 0)
		gtest.Assert(c.GetFloat64("name"), 0)
		gtest.Assert(c.GetFloats("name"), nil)
		gtest.Assert(c.GetInt("name"), 0)
		gtest.Assert(c.GetInt8("name"), 0)
		gtest.Assert(c.GetInt16("name"), 0)
		gtest.Assert(c.GetInt32("name"), 0)
		gtest.Assert(c.GetInt64("name"), 0)
		gtest.Assert(c.GetInts("name"), nil)
		gtest.Assert(c.GetUint("name"), 0)
		gtest.Assert(c.GetUint8("name"), 0)
		gtest.Assert(c.GetUint16("name"), 0)
		gtest.Assert(c.GetUint32("name"), 0)
		gtest.Assert(c.GetUint64("name"), 0)
		gtest.Assert(c.GetTime("name").Format("2006-01-02"), "0001-01-01")
		gtest.Assert(c.GetGTime("name"), nil)
		gtest.Assert(c.GetDuration("name").String(), "0s")
		name := struct {
			Name string
		}{}
		gtest.Assert(c.GetStruct("name", &name) == nil, false)

		c.Clear()

		arr, _ := gjson.Encode(
			g.Map{
				"name":   "gf",
				"time":   "2019-06-12",
				"person": g.Map{"name": "gf"},
				"floats": g.Slice{1, 2, 3},
			},
		)
		err = gfile.PutBytes(
			gfile.Join(configPath, "config.yml"),
			arr,
		)
		gtest.Assert(err, nil)
		gtest.Assert(c.GetTime("time").Format("2006-01-02"), "2019-06-12")
		gtest.Assert(c.GetGTime("time").Format("Y-m-d"), "2019-06-12")
		gtest.Assert(c.GetDuration("time").String(), "0s")

		err = c.GetStruct("person", &name)
		gtest.Assert(err, nil)
		gtest.Assert(name.Name, "gf")
		gtest.Assert(c.GetFloats("floats") == nil, false)
	})
}

func TestCfg_Instance(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gcfg.Instance("gf") != nil, true)
	})
	gtest.Case(t, func() {
		pwd := gfile.Pwd()
		gfile.Chdir(gfile.Join(gdebug.CallerDirectory(), "testdata"))
		defer gfile.Chdir(pwd)
		gtest.Assert(gcfg.Instance("c1") != nil, true)
		gtest.Assert(gcfg.Instance("c1").Get("my-config"), "1")
		gtest.Assert(gcfg.Instance("folder1/c1").Get("my-config"), "2")
	})
}

func TestCfg_Config(t *testing.T) {
	gtest.Case(t, func() {
		gcfg.SetContent("gf", "config.yml")
		gtest.Assert(gcfg.GetContent("config.yml"), "gf")
		gcfg.SetContent("gf1", "config.yml")
		gtest.Assert(gcfg.GetContent("config.yml"), "gf1")
		gcfg.RemoveContent("config.yml")
		gcfg.ClearContent()
		gtest.Assert(gcfg.GetContent("name"), "")
	})
}
