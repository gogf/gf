// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gi18n_test

import (
	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/gogf/gf/v2/os/gres/testdata/data"

	"context"
	"testing"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i18n := gi18n.New(gi18n.Options{
			Path: gtest.DataPath("i18n"),
		})
		i18n.SetLanguage("none")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")

		i18n.SetLanguage("ja")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "こんにちは世界")

		i18n.SetLanguage("zh-CN")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "你好世界")
		i18n.SetDelimiters("{$", "}")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")
		t.Assert(i18n.T(context.Background(), "{$hello}{$world}"), "你好世界")
	})

	gtest.C(t, func(t *gtest.T) {
		i18n := gi18n.New(gi18n.Options{
			Path: gtest.DataPath("i18n-file"),
		})
		i18n.SetLanguage("none")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")

		i18n.SetLanguage("ja")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "こんにちは世界")

		i18n.SetLanguage("zh-CN")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "你好世界")
	})

	gtest.C(t, func(t *gtest.T) {
		i18n := gi18n.New(gi18n.Options{
			Path: gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n-dir",
		})
		i18n.SetLanguage("none")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")

		i18n.SetLanguage("ja")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "こんにちは世界")

		i18n.SetLanguage("zh-CN")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "你好世界")
	})
}

func Test_TranslateFormat(t *testing.T) {
	// Tf
	gtest.C(t, func(t *gtest.T) {
		i18n := gi18n.New(gi18n.Options{
			Path: gtest.DataPath("i18n"),
		})
		i18n.SetLanguage("none")
		t.Assert(i18n.Tf(context.Background(), "{#hello}{#world} %d", 2020), "{#hello}{#world} 2020")

		i18n.SetLanguage("ja")
		t.Assert(i18n.Tf(context.Background(), "{#hello}{#world} %d", 2020), "こんにちは世界 2020")
	})
}

func Test_DefaultManager(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := gi18n.SetPath(gtest.DataPath("i18n"))
		t.AssertNil(err)

		gi18n.SetLanguage("none")
		t.Assert(gi18n.T(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")

		gi18n.SetLanguage("ja")
		t.Assert(gi18n.T(context.Background(), "{#hello}{#world}"), "こんにちは世界")

		gi18n.SetLanguage("zh-CN")
		t.Assert(gi18n.T(context.Background(), "{#hello}{#world}"), "你好世界")
	})

	gtest.C(t, func(t *gtest.T) {
		err := gi18n.SetPath(gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n-dir")
		t.AssertNil(err)

		gi18n.SetLanguage("none")
		t.Assert(gi18n.Translate(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")

		gi18n.SetLanguage("ja")
		t.Assert(gi18n.Translate(context.Background(), "{#hello}{#world}"), "こんにちは世界")

		gi18n.SetLanguage("zh-CN")
		t.Assert(gi18n.Translate(context.Background(), "{#hello}{#world}"), "你好世界")
	})
}

func Test_Instance(t *testing.T) {
	gres.Dump()
	gtest.C(t, func(t *gtest.T) {
		m := gi18n.Instance()
		err := m.SetPath("i18n-dir")
		t.AssertNil(err)
		m.SetLanguage("zh-CN")
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "你好世界")
	})

	gtest.C(t, func(t *gtest.T) {
		m := gi18n.Instance()
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "你好世界")
	})

	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.I18n().T(context.Background(), "{#hello}{#world}"), "你好世界")
	})
	// Default language is: en
	gtest.C(t, func(t *gtest.T) {
		m := gi18n.Instance(gconv.String(gtime.TimestampNano()))
		m.SetPath("i18n-dir")
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "HelloWorld")
	})
}

func Test_Resource(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.I18n("resource")
		err := m.SetPath("i18n-dir")
		t.AssertNil(err)

		m.SetLanguage("none")
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")

		m.SetLanguage("ja")
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "こんにちは世界")

		m.SetLanguage("zh-CN")
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "你好世界")
	})
}

func Test_SetCtxLanguage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		t.Assert(gi18n.LanguageFromCtx(ctx), "")
	})

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gi18n.LanguageFromCtx(nil), "")
	})

	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		ctx = gi18n.WithLanguage(ctx, "zh-CN")
		t.Assert(gi18n.LanguageFromCtx(ctx), "zh-CN")
	})

	gtest.C(t, func(t *gtest.T) {
		ctx := gi18n.WithLanguage(nil, "zh-CN")
		t.Assert(gi18n.LanguageFromCtx(ctx), "zh-CN")
	})

}
