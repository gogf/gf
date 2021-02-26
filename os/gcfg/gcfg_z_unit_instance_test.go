// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_Instance_Basic(t *testing.T) {
	config := `
array = [1.0, 2.0, 3.0]
v1 = 1.0
v2 = "true"
v3 = "off"
v4 = "1.234"

[redis]
  cache = "127.0.0.1:6379,1"
  disk = "127.0.0.1:6379,0"

`
	gtest.C(t, func(t *gtest.T) {
		path := DefaultConfigFile
		err := gfile.PutContents(path, config)
		t.Assert(err, nil)
		defer func() {
			t.Assert(gfile.Remove(path), nil)
		}()

		c := Instance()
		t.Assert(c.Get("v1"), 1)
		t.AssertEQ(c.GetInt("v1"), 1)
		t.AssertEQ(c.GetInt8("v1"), int8(1))
		t.AssertEQ(c.GetInt16("v1"), int16(1))
		t.AssertEQ(c.GetInt32("v1"), int32(1))
		t.AssertEQ(c.GetInt64("v1"), int64(1))
		t.AssertEQ(c.GetUint("v1"), uint(1))
		t.AssertEQ(c.GetUint8("v1"), uint8(1))
		t.AssertEQ(c.GetUint16("v1"), uint16(1))
		t.AssertEQ(c.GetUint32("v1"), uint32(1))
		t.AssertEQ(c.GetUint64("v1"), uint64(1))

		t.AssertEQ(c.GetVar("v1").String(), "1")
		t.AssertEQ(c.GetVar("v1").Bool(), true)
		t.AssertEQ(c.GetVar("v2").String(), "true")
		t.AssertEQ(c.GetVar("v2").Bool(), true)

		t.AssertEQ(c.GetString("v1"), "1")
		t.AssertEQ(c.GetFloat32("v4"), float32(1.234))
		t.AssertEQ(c.GetFloat64("v4"), float64(1.234))
		t.AssertEQ(c.GetString("v2"), "true")
		t.AssertEQ(c.GetBool("v2"), true)
		t.AssertEQ(c.GetBool("v3"), false)

		t.AssertEQ(c.Contains("v1"), true)
		t.AssertEQ(c.Contains("v2"), true)
		t.AssertEQ(c.Contains("v3"), true)
		t.AssertEQ(c.Contains("v4"), true)
		t.AssertEQ(c.Contains("v5"), false)

		t.AssertEQ(c.GetInts("array"), []int{1, 2, 3})
		t.AssertEQ(c.GetStrings("array"), []string{"1", "2", "3"})
		t.AssertEQ(c.GetArray("array"), []interface{}{1, 2, 3})
		t.AssertEQ(c.GetInterfaces("array"), []interface{}{1, 2, 3})
		t.AssertEQ(c.GetMap("redis"), map[string]interface{}{
			"disk":  "127.0.0.1:6379,0",
			"cache": "127.0.0.1:6379,1",
		})
		t.AssertEQ(c.FilePath(), gfile.Pwd()+gfile.Separator+path)
	})
}

func Test_Instance_AutoLocateConfigFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Instance("gf") != nil, true)
	})
	// Automatically locate the configuration file with supported file extensions.
	gtest.C(t, func(t *gtest.T) {
		pwd := gfile.Pwd()
		t.AssertNil(gfile.Chdir(gdebug.TestDataPath()))
		defer gfile.Chdir(pwd)
		t.Assert(Instance("c1") != nil, true)
		t.Assert(Instance("c1").Get("my-config"), "1")
		t.Assert(Instance("folder1/c1").Get("my-config"), "2")
	})
	// Automatically locate the configuration file with supported file extensions.
	gtest.C(t, func(t *gtest.T) {
		pwd := gfile.Pwd()
		t.AssertNil(gfile.Chdir(gdebug.TestDataPath("folder1")))
		defer gfile.Chdir(pwd)
		t.Assert(Instance("c2").Get("my-config"), 2)
	})
	// Default configuration file.
	gtest.C(t, func(t *gtest.T) {
		instances.Clear()
		pwd := gfile.Pwd()
		t.AssertNil(gfile.Chdir(gdebug.TestDataPath("default")))
		defer gfile.Chdir(pwd)
		t.Assert(Instance().Get("my-config"), 1)

		instances.Clear()
		t.AssertNil(genv.Set("GF_GCFG_FILE", "config.json"))
		defer genv.Set("GF_GCFG_FILE", "")
		t.Assert(Instance().Get("my-config"), 2)
	})
}

func Test_Instance_EnvPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		genv.Set("GF_GCFG_PATH", gdebug.TestDataPath("envpath"))
		defer genv.Set("GF_GCFG_PATH", "")
		t.Assert(Instance("c3") != nil, true)
		t.Assert(Instance("c3").Get("my-config"), "3")
		t.Assert(Instance("c4").Get("my-config"), "4")
		instances = gmap.NewStrAnyMap(true)
	})
}

func Test_Instance_EnvFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		genv.Set("GF_GCFG_PATH", gdebug.TestDataPath("envfile"))
		defer genv.Set("GF_GCFG_PATH", "")
		genv.Set("GF_GCFG_FILE", "c6.json")
		defer genv.Set("GF_GCFG_FILE", "")
		t.Assert(Instance().Get("my-config"), "6")
		instances = gmap.NewStrAnyMap(true)
	})
}
