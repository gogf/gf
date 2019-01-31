// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gstr_test

import (
    "gitee.com/johng/gf/g/string/gstr"
    "gitee.com/johng/gf/g/test/gtest"
    "testing"
)

func Test_Pos(t *testing.T) {
    gtest.Case(t, func() {
        s1 := "abcdEFGabcdefg"
        gtest.Assert(gstr.Pos(s1, "ab"),    0)
        gtest.Assert(gstr.Pos(s1, "ab", 2), 7)
    })
}

func Test_PosI(t *testing.T) {
    gtest.Case(t, func() {
        s1 := "abcdEFGabcdefg"
        gtest.Assert(gstr.PosI(s1, "zz"),   -1)
        gtest.Assert(gstr.PosI(s1, "ab"),    0)
        gtest.Assert(gstr.PosI(s1, "ef", 2), 4)
    })
}

func Test_PosR(t *testing.T) {
    gtest.Case(t, func() {
        s1 := "abcdEFGabcdefg"
        s2 := "abcdEFGz1cdeab"
        gtest.Assert(gstr.PosR(s1, "zz"),    -1)
        gtest.Assert(gstr.PosR(s1, "ab"),     7)
        gtest.Assert(gstr.PosR(s2, "ab", -2), 0)
        gtest.Assert(gstr.PosR(s1, "ef"),    11)
    })
}

func Test_PosRI(t *testing.T) {
    gtest.Case(t, func() {
        s1 := "abcdEFGabcdefg"
        s2 := "abcdEFGz1cdeab"
        gtest.Assert(gstr.PosRI(s1, "zz"),    -1)
        gtest.Assert(gstr.PosRI(s1, "AB"),     7)
        gtest.Assert(gstr.PosRI(s2, "AB", -2), 0)
        gtest.Assert(gstr.PosRI(s1, "EF"),    11)
    })
}