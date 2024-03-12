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
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestAdapterFile_Dump(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		t.Assert(c.GetFileName(), "config.yml")

		c.Dump()
		c.Data(ctx)
	})

	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("testdata/default/config.toml")
		t.AssertNil(err)

		c.Dump()
		c.Data(ctx)
		c.GetPaths()
	})

}
func TestAdapterFile_Available(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("testdata/default/config.toml")
		t.AssertNil(err)
		c.Available(ctx)
	})
}

func TestAdapterFile_SetPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		err = c.SetPath("/tmp")
		t.AssertNil(err)

		err = c.SetPath("notexist")
		t.AssertNE(err, nil)

		err = c.SetPath("testdata/c1.toml")
		t.AssertNE(err, nil)

		err = c.SetPath("")
		t.AssertNil(err)

		err = c.SetPath("gcfg.go")
		t.AssertNE(err, nil)

		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
	})
}

func TestAdapterFile_AddPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		err = c.AddPath("/tmp")
		t.AssertNil(err)

		err = c.AddPath("notexist")
		t.AssertNE(err, nil)

		err = c.SetPath("testdata/c1.toml")
		t.AssertNE(err, nil)

		err = c.SetPath("")
		t.AssertNil(err)

		err = c.AddPath("gcfg.go")
		t.AssertNE(err, nil)

		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
	})
}

func TestAdapterFile_SetViolenceCheck(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)
		c.SetViolenceCheck(true)
		v, err := c.Get(ctx, "name")
		t.AssertNE(err, nil)
		t.Assert(v, nil)
	})
}

func TestAdapterFile_FilePath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("config.yml")
		t.AssertNil(err)

		path, _ := c.GetFilePath("tmp")
		t.Assert(path, "")

		path, _ = c.GetFilePath("tmp")
		t.Assert(path, "")
	})
}

func TestAdapterFile_Content(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile()
		t.AssertNil(err)

		c.SetContent("gf", "config.yml")
		t.Assert(c.GetContent("config.yml"), "gf")
		c.SetContent("gf1", "config.yml")
		t.Assert(c.GetContent("config.yml"), "gf1")
		c.RemoveContent("config.yml")
		c.ClearContent()
		t.Assert(c.GetContent("name"), "")
	})
}

func TestAdapterFile_With_UTF8_BOM(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gcfg.NewAdapterFile("test-cfg-with-utf8-bom")
		t.AssertNil(err)

		t.Assert(c.SetPath("testdata"), nil)
		c.SetFileName("cfg-with-utf8-bom.toml")
		t.Assert(c.MustGet(ctx, "test.testInt"), 1)
		t.Assert(c.MustGet(ctx, "test.testStr"), "test")
	})
}

func TestAdapterFile_Set(t *testing.T) {
	config := `log-path = "logs"`
	gtest.C(t, func(t *gtest.T) {
		var (
			path = gcfg.DefaultConfigFileName
			err  = gfile.PutContents(path, config)
		)
		t.AssertNil(err)
		defer gfile.Remove(path)

		c, err := gcfg.New()
		t.Assert(c.MustGet(ctx, "log-path").String(), "logs")

		err = c.GetAdapter().(*gcfg.AdapterFile).Set("log-path", "custom-logs")
		t.AssertNil(err)
		t.Assert(c.MustGet(ctx, "log-path").String(), "custom-logs")
	})
}
