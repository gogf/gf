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
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Trim(" 123456\n "), "123456")
		t.Assert(gstr.Trim("#123456#;", "#;"), "123456")
	})
}

func Test_TrimStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimStr("gogo我爱gogo", "go"), "我爱")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimStr("gogo我爱gogo", "go", 1), "go我爱go")
		t.Assert(gstr.TrimStr("gogo我爱gogo", "go", 2), "我爱")
		t.Assert(gstr.TrimStr("gogo我爱gogo", "go", -1), "我爱")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimStr("啊我爱中国人啊", "啊"), "我爱中国人")
	})
}

func Test_TrimRight(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimRight(" 123456\n "), " 123456")
		t.Assert(gstr.TrimRight("#123456#;", "#;"), "#123456")
	})
}

func Test_TrimRightStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimRightStr("gogo我爱gogo", "go"), "gogo我爱")
		t.Assert(gstr.TrimRightStr("gogo我爱gogo", "go我爱gogo"), "go")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimRightStr("gogo我爱gogo", "go", 1), "gogo我爱go")
		t.Assert(gstr.TrimRightStr("gogo我爱gogo", "go", 2), "gogo我爱")
		t.Assert(gstr.TrimRightStr("gogo我爱gogo", "go", -1), "gogo我爱")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimRightStr("我爱中国人", "人"), "我爱中国")
		t.Assert(gstr.TrimRightStr("我爱中国人", "爱中国人"), "我")
	})
}

func Test_TrimLeft(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimLeft(" \r123456\n "), "123456\n ")
		t.Assert(gstr.TrimLeft("#;123456#;", "#;"), "123456#;")
	})
}

func Test_TrimLeftStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimLeftStr("gogo我爱gogo", "go"), "我爱gogo")
		t.Assert(gstr.TrimLeftStr("gogo我爱gogo", "gogo我爱go"), "go")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimLeftStr("gogo我爱gogo", "go", 1), "go我爱gogo")
		t.Assert(gstr.TrimLeftStr("gogo我爱gogo", "go", 2), "我爱gogo")
		t.Assert(gstr.TrimLeftStr("gogo我爱gogo", "go", -1), "我爱gogo")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.TrimLeftStr("我爱中国人", "我爱"), "中国人")
		t.Assert(gstr.TrimLeftStr("我爱中国人", "我爱中国"), "人")
	})
}
