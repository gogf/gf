// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcfg

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

var (
	ctx = context.TODO()
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
		t.Assert(c.MustGet(ctx, "v1"), 1)
		filepath, _ := c.GetAdapter().(*AdapterFile).GetFilePath()
		t.AssertEQ(filepath, gfile.Pwd()+gfile.Separator+path)
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
		t.Assert(Instance("c1").MustGet(ctx, "my-config"), "1")
		t.Assert(Instance("folder1/c1").MustGet(ctx, "my-config"), "2")
	})
	// Automatically locate the configuration file with supported file extensions.
	gtest.C(t, func(t *gtest.T) {
		pwd := gfile.Pwd()
		t.AssertNil(gfile.Chdir(gdebug.TestDataPath("folder1")))
		defer gfile.Chdir(pwd)
		t.Assert(Instance("c2").MustGet(ctx, "my-config"), 2)
	})
	// Default configuration file.
	gtest.C(t, func(t *gtest.T) {
		localInstances.Clear()
		pwd := gfile.Pwd()
		t.AssertNil(gfile.Chdir(gdebug.TestDataPath("default")))
		defer gfile.Chdir(pwd)
		t.Assert(Instance().MustGet(ctx, "my-config"), 1)

		localInstances.Clear()
		t.AssertNil(genv.Set("GF_GCFG_FILE", "config.json"))
		defer genv.Set("GF_GCFG_FILE", "")
		t.Assert(Instance().MustGet(ctx, "my-config"), 2)
	})
}

func Test_Instance_EnvPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		genv.Set("GF_GCFG_PATH", gdebug.TestDataPath("envpath"))
		defer genv.Set("GF_GCFG_PATH", "")
		t.Assert(Instance("c3") != nil, true)
		t.Assert(Instance("c3").MustGet(ctx, "my-config"), "3")
		t.Assert(Instance("c4").MustGet(ctx, "my-config"), "4")
		localInstances = gmap.NewStrAnyMap(true)
	})
}

func Test_Instance_EnvFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		genv.Set("GF_GCFG_PATH", gdebug.TestDataPath("envfile"))
		defer genv.Set("GF_GCFG_PATH", "")
		genv.Set("GF_GCFG_FILE", "c6.json")
		defer genv.Set("GF_GCFG_FILE", "")
		t.Assert(Instance().MustGet(ctx, "my-config"), "6")
		localInstances = gmap.NewStrAnyMap(true)
	})
}
