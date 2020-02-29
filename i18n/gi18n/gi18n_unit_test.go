// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gi18n_test

import (
	"testing"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/i18n/gi18n"

	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/gfile"

	"github.com/gogf/gf/test/gtest"

	_ "github.com/gogf/gf/os/gres/testdata/data"
)

func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		t := gi18n.New(gi18n.Options{
			Path: gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n",
		})
		t.SetLanguage("none")
		gtest.Assert(t.T("{#hello}{#world}"), "{#hello}{#world}")

		t.SetLanguage("ja")
		gtest.Assert(t.T("{#hello}{#world}"), "こんにちは世界")

		t.SetLanguage("zh-CN")
		gtest.Assert(t.T("{#hello}{#world}"), "你好世界")
		t.SetDelimiters("{$", "}")
		gtest.Assert(t.T("{#hello}{#world}"), "{#hello}{#world}")
		gtest.Assert(t.T("{$hello}{$world}"), "你好世界")
	})

	gtest.Case(t, func() {
		t := gi18n.New(gi18n.Options{
			Path: gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n-file",
		})
		t.SetLanguage("none")
		gtest.Assert(t.T("{#hello}{#world}"), "{#hello}{#world}")

		t.SetLanguage("ja")
		gtest.Assert(t.T("{#hello}{#world}"), "こんにちは世界")

		t.SetLanguage("zh-CN")
		gtest.Assert(t.T("{#hello}{#world}"), "你好世界")
	})

	gtest.Case(t, func() {
		t := gi18n.New(gi18n.Options{
			Path: gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n-dir",
		})
		t.SetLanguage("none")
		gtest.Assert(t.T("{#hello}{#world}"), "{#hello}{#world}")

		t.SetLanguage("ja")
		gtest.Assert(t.T("{#hello}{#world}"), "こんにちは世界")

		t.SetLanguage("zh-CN")
		gtest.Assert(t.T("{#hello}{#world}"), "你好世界")
	})
}

func Test_DefaultManager(t *testing.T) {
	gtest.Case(t, func() {
		err := gi18n.SetPath(gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n")
		gtest.Assert(err, nil)

		gi18n.SetLanguage("none")
		gtest.Assert(gi18n.T("{#hello}{#world}"), "{#hello}{#world}")

		gi18n.SetLanguage("ja")
		gtest.Assert(gi18n.T("{#hello}{#world}"), "こんにちは世界")

		gi18n.SetLanguage("zh-CN")
		gtest.Assert(gi18n.T("{#hello}{#world}"), "你好世界")
	})

	gtest.Case(t, func() {
		err := gi18n.SetPath(gdebug.CallerDirectory() + gfile.Separator + "testdata" + gfile.Separator + "i18n-dir")
		gtest.Assert(err, nil)

		gi18n.SetLanguage("none")
		gtest.Assert(gi18n.Translate("{#hello}{#world}"), "{#hello}{#world}")

		gi18n.SetLanguage("ja")
		gtest.Assert(gi18n.Translate("{#hello}{#world}"), "こんにちは世界")

		gi18n.SetLanguage("zh-CN")
		gtest.Assert(gi18n.Translate("{#hello}{#world}"), "你好世界")
	})
}

func Test_Instance(t *testing.T) {
	gres.Dump()
	gtest.Case(t, func() {
		m := gi18n.Instance()
		err := m.SetPath("i18n-dir")
		gtest.Assert(err, nil)
		m.SetLanguage("zh-CN")
		gtest.Assert(m.T("{#hello}{#world}"), "你好世界")
	})

	gtest.Case(t, func() {
		m := gi18n.Instance()
		gtest.Assert(m.T("{#hello}{#world}"), "你好世界")
	})

	gtest.Case(t, func() {
		gtest.Assert(g.I18n().T("{#hello}{#world}"), "你好世界")
	})

	gtest.Case(t, func() {
		m := gi18n.Instance(gconv.String(gtime.TimestampNano()))
		gtest.Assert(m.T("{#hello}{#world}"), "{#hello}{#world}")
	})
}

func Test_Resource(t *testing.T) {
	gtest.Case(t, func() {
		m := g.I18n("resource")
		err := m.SetPath("i18n-dir")
		gtest.Assert(err, nil)

		m.SetLanguage("none")
		gtest.Assert(m.T("{#hello}{#world}"), "{#hello}{#world}")

		m.SetLanguage("ja")
		gtest.Assert(m.T("{#hello}{#world}"), "こんにちは世界")

		m.SetLanguage("zh-CN")
		gtest.Assert(m.T("{#hello}{#world}"), "你好世界")
	})
}
