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

func Test_Pos(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(gstr.Pos(s1, "ab"), 0)
		t.Assert(gstr.Pos(s1, "ab", 2), 7)
		t.Assert(gstr.Pos(s1, "abd", 0), -1)
		t.Assert(gstr.Pos(s1, "e", -4), 11)
	})
}

func Test_PosI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(gstr.PosI(s1, "zz"), -1)
		t.Assert(gstr.PosI(s1, "ab"), 0)
		t.Assert(gstr.PosI(s1, "ef", 2), 4)
		t.Assert(gstr.PosI(s1, "abd", 0), -1)
		t.Assert(gstr.PosI(s1, "E", -4), 11)
	})
}

func Test_PosR(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(gstr.PosR(s1, "zz"), -1)
		t.Assert(gstr.PosR(s1, "ab"), 7)
		t.Assert(gstr.PosR(s2, "ab", -2), 0)
		t.Assert(gstr.PosR(s1, "ef"), 11)
		t.Assert(gstr.PosR(s1, "abd", 0), -1)
		t.Assert(gstr.PosR(s1, "e", -4), -1)
	})
}

func Test_PosRI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(gstr.PosRI(s1, "zz"), -1)
		t.Assert(gstr.PosRI(s1, "AB"), 7)
		t.Assert(gstr.PosRI(s2, "AB", -2), 0)
		t.Assert(gstr.PosRI(s1, "EF"), 11)
		t.Assert(gstr.PosRI(s1, "abd", 0), -1)
		t.Assert(gstr.PosRI(s1, "e", -5), 4)
	})
}
