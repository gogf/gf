// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gi18n_test

import (
	"testing"

	"github.com/gogf/gf/util/gi18n"

	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/gfile"

	"github.com/gogf/gf/test/gtest"
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
	})
}
