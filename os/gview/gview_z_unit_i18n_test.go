// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_I18n(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "你好世界!"`
		expect2 := `john says "こんにちは世界!"`
		expect3 := `john says "{#hello}{#world}!"`

		g.I18n().SetPath(gtest.DataPath("i18n"))

		g.I18n().SetLanguage("zh-CN")
		result1, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)

		g.I18n().SetLanguage("ja")
		result2, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result2, expect2)

		g.I18n().SetLanguage("none")
		result3, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result3, expect3)
	})
	gtest.C(t, func(t *gtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "你好世界!"`
		expect2 := `john says "こんにちは世界!"`
		expect3 := `john says "{#hello}{#world}!"`

		g.I18n().SetPath(gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n")

		result1, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "zh-CN",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)

		result2, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "ja",
		})
		t.AssertNil(err)
		t.Assert(result2, expect2)

		result3, err := g.View().ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "none",
		})
		t.AssertNil(err)
		t.Assert(result3, expect3)
	})
	// gi18n manager is nil
	gtest.C(t, func(t *gtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "{#hello}{#world}!"`

		g.I18n().SetPath(gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n")

		view := gview.New()
		view.SetI18n(nil)
		result1, err := view.ParseContent(context.TODO(), content, g.Map{
			"name":         "john",
			"I18nLanguage": "zh-CN",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)
	})
	// SetLanguage in context
	gtest.C(t, func(t *gtest.T) {
		content := `{{.name}} says "{#hello}{#world}!"`
		expect1 := `john says "你好世界!"`
		ctx := gctx.New()
		g.I18n().SetPath(gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n")
		ctx = gi18n.WithLanguage(ctx, "zh-CN")
		t.Log(gi18n.LanguageFromCtx(ctx))

		view := gview.New()

		result1, err := view.ParseContent(ctx, content, g.Map{
			"name": "john",
		})
		t.AssertNil(err)
		t.Assert(result1, expect1)
	})

}
