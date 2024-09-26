// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gi18n_test

import (
	"time"

	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/os/gctx"

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
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "{#hello}{#world}")
		t.Assert(i18n.T(context.Background(), "{$你好} {$世界}"), "hello world")
		// undefined variables.
		t.Assert(i18n.T(context.Background(), "{$你好1}{$世界1}"), "{$你好1}{$世界1}")
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
		t.Assert(i18n.T(context.Background(), "{#你好} {#世界}"), "hello world")
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
	gtest.C(t, func(t *gtest.T) {
		m := gi18n.Instance()
		err := m.SetPath(gtest.DataPath("i18n-dir"))
		t.AssertNil(err)
		m.SetLanguage("zh-CN")
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "你好世界")
		t.Assert(m.T(context.Background(), "{#你好} {#世界}"), "hello world")
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
		m.SetPath(gtest.DataPath("i18n-dir"))
		t.Assert(m.T(context.Background(), "{#hello}{#world}"), "HelloWorld")
	})
}

func Test_Resource(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.I18n("resource")
		err := m.SetPath(gtest.DataPath("i18n-dir"))
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
		ctx := gi18n.WithLanguage(context.Background(), "zh-CN")
		t.Assert(gi18n.LanguageFromCtx(ctx), "zh-CN")
	})

}

func Test_GetContent(t *testing.T) {
	i18n := gi18n.New(gi18n.Options{
		Path: gtest.DataPath("i18n-file"),
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(i18n.GetContent(context.Background(), "hello"), "Hello")

		ctx := gi18n.WithLanguage(context.Background(), "zh-CN")
		t.Assert(i18n.GetContent(ctx, "hello"), "你好")

		ctx = gi18n.WithLanguage(context.Background(), "unknown")
		t.Assert(i18n.GetContent(ctx, "hello"), "")
	})
}

func Test_PathInResource(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		binContent, err := gres.Pack(gtest.DataPath("i18n"))
		t.AssertNil(err)
		err = gres.Add(gbase64.EncodeToString(binContent))
		t.AssertNil(err)

		i18n := gi18n.New()
		i18n.SetLanguage("zh-CN")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "你好世界")

		err = i18n.SetPath("i18n")
		t.AssertNil(err)
		i18n.SetLanguage("ja")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "こんにちは世界")
	})
}

func Test_PathInNormal(t *testing.T) {
	// Copy i18n files to current directory.
	gfile.CopyDir(gtest.DataPath("i18n"), gfile.Join(gdebug.CallerDirectory(), "manifest/i18n"))
	// Remove copied files after testing.
	defer gfile.Remove(gfile.Join(gdebug.CallerDirectory(), "manifest"))

	i18n := gi18n.New()

	gtest.C(t, func(t *gtest.T) {
		i18n.SetLanguage("zh-CN")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "你好世界")
		// Set not exist path.
		err := i18n.SetPath("i18n-not-exist")
		t.AssertNE(err, nil)
		err = i18n.SetPath("")
		t.AssertNE(err, nil)
		i18n.SetLanguage("ja")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}"), "こんにちは世界")
	})

	// Change language file content.
	gtest.C(t, func(t *gtest.T) {
		i18n.SetLanguage("en")
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}{#name}"), "HelloWorld{#name}")
		err := gfile.PutContentsAppend(gfile.Join(gdebug.CallerDirectory(), "manifest/i18n/en.toml"), "\nname = \"GoFrame\"")
		t.AssertNil(err)
		// Wait for the file modification time to change.
		time.Sleep(10 * time.Millisecond)
		t.Assert(i18n.T(context.Background(), "{#hello}{#world}{#name}"), "HelloWorldGoFrame")
	})

	// Add new language
	gtest.C(t, func(t *gtest.T) {
		err := gfile.PutContents(gfile.Join(gdebug.CallerDirectory(), "manifest/i18n/en-US.toml"), "lang = \"en-US\"")
		t.AssertNil(err)
		// Wait for the file modification time to change.
		time.Sleep(10 * time.Millisecond)
		i18n.SetLanguage("en-US")
		t.Assert(i18n.T(context.Background(), "{#lang}"), "en-US")
	})
}
