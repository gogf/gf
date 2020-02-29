// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/gcfg"
	"testing"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_View(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertNE(View(), nil)
		b, e := View().ParseContent(`{{"我是中国人" | substr 2 -1}}`, nil)
		gtest.Assert(e, nil)
		gtest.Assert(b, "中国人")
	})
	gtest.Case(t, func() {
		tpl := "t.tpl"
		err := gfile.PutContents(tpl, `{{"我是中国人" | substr 2 -1}}`)
		gtest.Assert(err, nil)
		defer gfile.Remove(tpl)

		b, e := View().Parse("t.tpl", nil)
		gtest.Assert(e, nil)
		gtest.Assert(b, "中国人")
	})
	gtest.Case(t, func() {
		path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.TimestampNano())
		tpl := fmt.Sprintf(`%s/%s`, path, "t.tpl")
		err := gfile.PutContents(tpl, `{{"我是中国人" | substr 2 -1}}`)
		gtest.Assert(err, nil)
		defer gfile.Remove(tpl)
		err = View().AddPath(path)
		gtest.Assert(err, nil)

		b, e := View().Parse("t.tpl", nil)
		gtest.Assert(e, nil)
		gtest.Assert(b, "中国人")
	})
}

func Test_View_Config(t *testing.T) {
	// view1 test1
	gtest.Case(t, func() {
		dirPath := gfile.Join(gdebug.CallerDirectory(), "testdata", "view1")
		gcfg.SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer gcfg.ClearContent()
		defer instances.Clear()

		view := View("test1")
		gtest.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		gtest.Assert(err, nil)

		str := `hello ${.name},version:${.version}`
		view.Assigns(map[string]interface{}{"version": "1.9.0"})
		result, err := view.ParseContent(str, nil)
		gtest.Assert(err, nil)
		gtest.Assert(result, "hello test1,version:1.9.0")

		result, err = view.ParseDefault()
		gtest.Assert(err, nil)
		gtest.Assert(result, "test1:test1")
	})
	// view1 test2
	gtest.Case(t, func() {
		dirPath := gfile.Join(gdebug.CallerDirectory(), "testdata", "view1")
		gcfg.SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer gcfg.ClearContent()
		defer instances.Clear()

		view := View("test2")
		gtest.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		gtest.Assert(err, nil)

		str := `hello #{.name},version:#{.version}`
		view.Assigns(map[string]interface{}{"version": "1.9.0"})
		result, err := view.ParseContent(str, nil)
		gtest.Assert(err, nil)
		gtest.Assert(result, "hello test2,version:1.9.0")

		result, err = view.ParseDefault()
		gtest.Assert(err, nil)
		gtest.Assert(result, "test2:test2")
	})
	// view2
	gtest.Case(t, func() {
		dirPath := gfile.Join(gdebug.CallerDirectory(), "testdata", "view2")
		gcfg.SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer gcfg.ClearContent()
		defer instances.Clear()

		view := View()
		gtest.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		gtest.Assert(err, nil)

		str := `hello {.name},version:{.version}`
		view.Assigns(map[string]interface{}{"version": "1.9.0"})
		result, err := view.ParseContent(str, nil)
		gtest.Assert(err, nil)
		gtest.Assert(result, "hello test,version:1.9.0")

		result, err = view.ParseDefault()
		gtest.Assert(err, nil)
		gtest.Assert(result, "test:test")
	})
	// view2
	gtest.Case(t, func() {
		dirPath := gfile.Join(gdebug.CallerDirectory(), "testdata", "view2")
		gcfg.SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer gcfg.ClearContent()
		defer instances.Clear()

		view := View("test100")
		gtest.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		gtest.Assert(err, nil)

		str := `hello {.name},version:{.version}`
		view.Assigns(map[string]interface{}{"version": "1.9.0"})
		result, err := view.ParseContent(str, nil)
		gtest.Assert(err, nil)
		gtest.Assert(result, "hello test,version:1.9.0")

		result, err = view.ParseDefault()
		gtest.Assert(err, nil)
		gtest.Assert(result, "test:test")
	})
}
