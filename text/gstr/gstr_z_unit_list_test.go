// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_List2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.List2("1:2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.List2("1:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.List2("1", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.List2("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.List2("1:2:3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2:3")
	})
}

func Test_ListAndTrim2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.ListAndTrim2("1::2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.ListAndTrim2("1::", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.ListAndTrim2("1:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.ListAndTrim2("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2 := gstr.ListAndTrim2("1::2::3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2:3")
	})
}

func Test_List3(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.List3("1:2:3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "3")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.List3("1:2:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.List3("1:2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.List3("1:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.List3("1", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.List3("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.List3("1:2:3:4", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "3:4")
	})
}

func Test_ListAndTrim3(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.ListAndTrim3("1::2:3", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "3")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.ListAndTrim3("1::2:", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.ListAndTrim3("1::2", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "2")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.ListAndTrim3("1::", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.ListAndTrim3("1::", ":")
		t.Assert(p1, "1")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
	gtest.C(t, func(t *gtest.T) {
		p1, p2, p3 := gstr.ListAndTrim3("", ":")
		t.Assert(p1, "")
		t.Assert(p2, "")
		t.Assert(p3, "")
	})
}
