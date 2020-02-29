// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package grand_test

import (
	"github.com/gogf/gf/text/gstr"
	"testing"

	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/grand"
)

func Test_Intn(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 1000000; i++ {
			n := grand.Intn(100)
			gtest.AssertLT(n, 100)
			gtest.AssertGE(n, 0)
		}
		for i := 0; i < 1000000; i++ {
			n := grand.Intn(-100)
			gtest.AssertLE(n, 0)
			gtest.Assert(n, -100)
		}
	})
}

func Test_Meet(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.Meet(100, 100), true)
		}
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.Meet(0, 100), false)
		}
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.Meet(50, 100), []bool{true, false})
		}
	})
}

func Test_MeetProb(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.MeetProb(1), true)
		}
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.MeetProb(0), false)
		}
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.MeetProb(0.5), []bool{true, false})
		}
	})
}

func Test_N(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.N(1, 2), []int{1, 2})
		}
	})
}

func Test_Rand(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.N(1, 2), []int{1, 2})
		}
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.N(-1, 2), []int{-1, 0, 1, 2})
		}
	})
}

func Test_Str(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.S(5)), 5)
		}
	})
}

func Test_RandS(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.S(5)), 5)
		}
	})
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.S(5, true)), 5)
		}
	})
}

func Test_RandStr(t *testing.T) {
	str := "我爱GoFrame"
	gtest.Case(t, func() {
		for i := 0; i < 10; i++ {
			s := grand.Str(str, 100000)
			gtest.Assert(gstr.Contains(s, "我"), true)
			gtest.Assert(gstr.Contains(s, "爱"), true)
			gtest.Assert(gstr.Contains(s, "G"), true)
			gtest.Assert(gstr.Contains(s, "o"), true)
			gtest.Assert(gstr.Contains(s, "F"), true)
			gtest.Assert(gstr.Contains(s, "r"), true)
			gtest.Assert(gstr.Contains(s, "a"), true)
			gtest.Assert(gstr.Contains(s, "m"), true)
			gtest.Assert(gstr.Contains(s, "e"), true)
			gtest.Assert(gstr.Contains(s, "w"), false)
		}
	})
}

func Test_Digits(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.Digits(5)), 5)
		}
	})
}

func Test_RandDigits(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.Digits(5)), 5)
		}
	})
}

func Test_Letters(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.Letters(5)), 5)
		}
	})
}

func Test_RandLetters(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.Letters(5)), 5)
		}
	})
}

func Test_Perm(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.Perm(5), []int{0, 1, 2, 3, 4})
		}
	})
}
