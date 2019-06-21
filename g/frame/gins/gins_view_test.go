// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/g/frame/gins"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_View(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertNE(gins.View(), nil)
		b, e := gins.View().ParseContent(`{{"我是中国人" | substr 2 -1}}`, nil)
		gtest.Assert(e, nil)
		gtest.Assert(string(b), "中国人")
	})
	gtest.Case(t, func() {
		tpl := "t.tpl"
		err := gfile.PutContents(tpl, `{{"我是中国人" | substr 2 -1}}`)
		gtest.Assert(err, nil)
		defer gfile.Remove(tpl)

		b, e := gins.View().Parse("t.tpl", nil)
		gtest.Assert(e, nil)
		gtest.Assert(string(b), "中国人")
	})
	gtest.Case(t, func() {
		path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.Nanosecond())
		tpl := fmt.Sprintf(`%s/%s`, path, "t.tpl")
		err := gfile.PutContents(tpl, `{{"我是中国人" | substr 2 -1}}`)
		gtest.Assert(err, nil)
		defer gfile.Remove(tpl)
		err = gins.View().AddPath(path)
		gtest.Assert(err, nil)

		b, e := gins.View().Parse("t.tpl", nil)
		gtest.Assert(e, nil)
		gtest.Assert(string(b), "中国人")
	})
}
