// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
)

func Test_Trim(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.Trim(" 123456\n "), "123456")
		gtest.Assert(gstr.Trim("#123456#;", "#;"), "123456")
	})
}

func Test_TrimStr(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimStr("gogo我爱gogo", "go"), "我爱")
	})
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimStr("啊我爱中国人啊", "啊"), "我爱中国人")
	})
}

func Test_TrimRight(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimRight(" 123456\n "), " 123456")
		gtest.Assert(gstr.TrimRight("#123456#;", "#;"), "#123456")
	})
}

func Test_TrimRightStr(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimRightStr("gogo我爱gogo", "go"), "gogo我爱")
		gtest.Assert(gstr.TrimRightStr("gogo我爱gogo", "go我爱gogo"), "go")
	})
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimRightStr("我爱中国人", "人"), "我爱中国")
		gtest.Assert(gstr.TrimRightStr("我爱中国人", "爱中国人"), "我")
	})
}

func Test_TrimLeft(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimLeft(" \r123456\n "), "123456\n ")
		gtest.Assert(gstr.TrimLeft("#;123456#;", "#;"), "123456#;")
	})
}

func Test_TrimLeftStr(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimLeftStr("gogo我爱gogo", "go"), "我爱gogo")
		gtest.Assert(gstr.TrimLeftStr("gogo我爱gogo", "gogo我爱go"), "go")
	})
	gtest.Case(t, func() {
		gtest.Assert(gstr.TrimLeftStr("我爱中国人", "我爱"), "中国人")
		gtest.Assert(gstr.TrimLeftStr("我爱中国人", "我爱中国"), "人")
	})
}
