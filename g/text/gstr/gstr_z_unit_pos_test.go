// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/text/gstr"
	"testing"
)

func Test_Pos(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFGabcdefg"
		gtest.Assert(gstr.Pos(s1, "ab"), 0)
		gtest.Assert(gstr.Pos(s1, "ab", 2), 7)
		gtest.Assert(gstr.Pos(s1, "abd", 0), -1)
		gtest.Assert(gstr.Pos(s1, "e", -4), 11)
	})
}

func Test_PosI(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFGabcdefg"
		gtest.Assert(gstr.PosI(s1, "zz"), -1)
		gtest.Assert(gstr.PosI(s1, "ab"), 0)
		gtest.Assert(gstr.PosI(s1, "ef", 2), 4)
		gtest.Assert(gstr.PosI(s1, "abd", 0), -1)
		gtest.Assert(gstr.PosI(s1, "E", -4), 11)
	})
}

func Test_PosR(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		gtest.Assert(gstr.PosR(s1, "zz"), -1)
		gtest.Assert(gstr.PosR(s1, "ab"), 7)
		gtest.Assert(gstr.PosR(s2, "ab", -2), 0)
		gtest.Assert(gstr.PosR(s1, "ef"), 11)
		gtest.Assert(gstr.PosR(s1, "abd", 0), -1)
		gtest.Assert(gstr.PosR(s1, "e", -4), -1)
	})
}

func Test_PosRI(t *testing.T) {
	gtest.Case(t, func() {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		gtest.Assert(gstr.PosRI(s1, "zz"), -1)
		gtest.Assert(gstr.PosRI(s1, "AB"), 7)
		gtest.Assert(gstr.PosRI(s2, "AB", -2), 0)
		gtest.Assert(gstr.PosRI(s1, "EF"), 11)
		gtest.Assert(gstr.PosRI(s1, "abd", 0), -1)
		gtest.Assert(gstr.PosRI(s1, "e", -5), 4)
	})
}
