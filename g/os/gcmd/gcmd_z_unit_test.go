// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd

import (
	"os"
	"testing"

	"github.com/gogf/gf/g/test/gtest"
)

func Test_ValueAndOption(t *testing.T) {
	os.Args = []string{"v1", "v2", "--o1=111", "-o2=222"}
	doInit()
	gtest.Case(t, func() {
		gtest.Assert(Value.GetAll(), []string{"v1", "v2"})
		gtest.Assert(Value.Get(0), "v1")
		gtest.Assert(Value.Get(1), "v2")
		gtest.Assert(Value.Get(2), "")
		gtest.Assert(Value.Get(2, "1"), "1")
		gtest.Assert(Value.GetVar(1, "1").String(), "v2")
		gtest.Assert(Value.GetVar(2, "1").String(), "1")

		gtest.Assert(Option.GetAll(), map[string]string{"o1": "111", "o2": "222"})
		gtest.Assert(Option.Get("o1"), "111")
		gtest.Assert(Option.Get("o2"), "222")
		gtest.Assert(Option.Get("o3", "1"), "1")
		gtest.Assert(Option.GetVar("o2", "1").String(), "222")
		gtest.Assert(Option.GetVar("o3", "1").String(), "1")

	})
}

func Test_Handle(t *testing.T) {
	os.Args = []string{"gf", "gf"}
	doInit()
	gtest.Case(t, func() {
		num := 1
		BindHandle("gf", func() {
			num += 1
		})
		RunHandle("gf")
		gtest.AssertEQ(num, 2)
		AutoRun()
		gtest.AssertEQ(num, 3)
	})
}
