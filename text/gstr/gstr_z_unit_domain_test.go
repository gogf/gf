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

func Test_IsSubDomain(t *testing.T) {
	gtest.Case(t, func() {
		main := "goframe.org"
		gtest.Assert(gstr.IsSubDomain("goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("s.goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("s.s.goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("johng.cn", main), false)
		gtest.Assert(gstr.IsSubDomain("s.johng.cn", main), false)
		gtest.Assert(gstr.IsSubDomain("s.s.johng.cn", main), false)
	})
	gtest.Case(t, func() {
		main := "*.goframe.org"
		gtest.Assert(gstr.IsSubDomain("goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("s.goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("s.s.goframe.org", main), false)
		gtest.Assert(gstr.IsSubDomain("johng.cn", main), false)
		gtest.Assert(gstr.IsSubDomain("s.johng.cn", main), false)
		gtest.Assert(gstr.IsSubDomain("s.s.johng.cn", main), false)
	})
	gtest.Case(t, func() {
		main := "*.*.goframe.org"
		gtest.Assert(gstr.IsSubDomain("goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("s.goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("s.s.goframe.org", main), true)
		gtest.Assert(gstr.IsSubDomain("s.s.s.goframe.org", main), false)
		gtest.Assert(gstr.IsSubDomain("johng.cn", main), false)
		gtest.Assert(gstr.IsSubDomain("s.johng.cn", main), false)
		gtest.Assert(gstr.IsSubDomain("s.s.johng.cn", main), false)
	})
}
