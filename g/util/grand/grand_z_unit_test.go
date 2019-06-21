// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package grand_test

import (
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/grand"
	"testing"
)

func Test_Intn(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 1000000; i++ {
			n := grand.Intn(100)
			gtest.AssertLT(n, 100)
			gtest.AssertGTE(n, 0)
		}
		for i := 0; i < 1000000; i++ {
			n := grand.Intn(-100)
			gtest.AssertLTE(n, 0)
			gtest.AssertGT(n, -100)
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
			gtest.Assert(grand.Rand(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			gtest.Assert(grand.Rand(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.Rand(1, 2), []int{1, 2})
		}
		for i := 0; i < 100; i++ {
			gtest.AssertIN(grand.Rand(-1, 2), []int{-1, 0, 1, 2})
		}
	})
}

func Test_Str(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.Str(5)), 5)
		}
	})
}

func Test_RandStr(t *testing.T) {
	gtest.Case(t, func() {
		for i := 0; i < 100; i++ {
			gtest.Assert(len(grand.RandStr(5)), 5)
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
			gtest.Assert(len(grand.RandDigits(5)), 5)
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
			gtest.Assert(len(grand.RandLetters(5)), 5)
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
