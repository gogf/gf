// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/internal/instance"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_View(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(View(), nil)
		b, e := View().ParseContent(context.TODO(), `{{"我是中国人" | substr 2 -1}}`, nil)
		t.Assert(e, nil)
		t.Assert(b, "中国")
	})
	gtest.C(t, func(t *gtest.T) {
		tpl := "t.tpl"
		err := gfile.PutContents(tpl, `{{"我是中国人" | substr 2 -1}}`)
		t.AssertNil(err)
		defer gfile.Remove(tpl)

		b, e := View().Parse(context.TODO(), "t.tpl", nil)
		t.Assert(e, nil)
		t.Assert(b, "中国")
	})
	gtest.C(t, func(t *gtest.T) {
		path := fmt.Sprintf(`%s/%d`, gfile.Temp(), gtime.TimestampNano())
		tpl := fmt.Sprintf(`%s/%s`, path, "t.tpl")
		err := gfile.PutContents(tpl, `{{"我是中国人" | substr 2 -1}}`)
		t.AssertNil(err)
		defer gfile.Remove(tpl)
		err = View().AddPath(path)
		t.AssertNil(err)

		b, e := View().Parse(context.TODO(), "t.tpl", nil)
		t.Assert(e, nil)
		t.Assert(b, "中国")
	})
}

func Test_View_Config(t *testing.T) {
	var ctx = context.TODO()
	// view1 test1
	gtest.C(t, func(t *gtest.T) {
		dirPath := gtest.DataPath("view1")
		Config().GetAdapter().(*gcfg.AdapterFile).SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer Config().GetAdapter().(*gcfg.AdapterFile).ClearContent()
		defer instance.Clear()

		view := View("test1")
		t.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		t.AssertNil(err)

		str := `hello ${.name},version:${.version}`
		view.Assigns(map[string]any{"version": "1.9.0"})
		result, err := view.ParseContent(ctx, str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello test1,version:1.9.0")

		result, err = view.ParseDefault(ctx)
		t.AssertNil(err)
		t.Assert(result, "test1:test1")
	})
	// view1 test2
	gtest.C(t, func(t *gtest.T) {
		dirPath := gtest.DataPath("view1")
		Config().GetAdapter().(*gcfg.AdapterFile).SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer Config().GetAdapter().(*gcfg.AdapterFile).ClearContent()
		defer instance.Clear()

		view := View("test2")
		t.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		t.AssertNil(err)

		str := `hello #{.name},version:#{.version}`
		view.Assigns(map[string]any{"version": "1.9.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello test2,version:1.9.0")

		result, err = view.ParseDefault(context.TODO())
		t.AssertNil(err)
		t.Assert(result, "test2:test2")
	})
	// view2
	gtest.C(t, func(t *gtest.T) {
		dirPath := gtest.DataPath("view2")
		Config().GetAdapter().(*gcfg.AdapterFile).SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer Config().GetAdapter().(*gcfg.AdapterFile).ClearContent()
		defer instance.Clear()

		view := View()
		t.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		t.AssertNil(err)

		str := `hello {.name},version:{.version}`
		view.Assigns(map[string]any{"version": "1.9.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello test,version:1.9.0")

		result, err = view.ParseDefault(context.TODO())
		t.AssertNil(err)
		t.Assert(result, "test:test")
	})
	// view2
	gtest.C(t, func(t *gtest.T) {
		dirPath := gtest.DataPath("view2")
		Config().GetAdapter().(*gcfg.AdapterFile).SetContent(gfile.GetContents(gfile.Join(dirPath, "config.toml")))
		defer Config().GetAdapter().(*gcfg.AdapterFile).ClearContent()
		defer instance.Clear()

		view := View("test100")
		t.AssertNE(view, nil)
		err := view.AddPath(dirPath)
		t.AssertNil(err)

		str := `hello {.name},version:{.version}`
		view.Assigns(map[string]any{"version": "1.9.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello test,version:1.9.0")

		result, err = view.ParseDefault(context.TODO())
		t.AssertNil(err)
		t.Assert(result, "test:test")
	})
}
