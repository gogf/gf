// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd

import (
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"testing"
)


func Test_ValueAndOption(t *testing.T) {
	os.Args = []string{"v1", "v2", "--o1=111", "-o2=222"}
	doInit()
	gtest.Case(t, func() {
		gtest.Assert(Value.GetAll(), []string{"v1", "v2"})
		gtest.Assert(Value.Get(0), "v1")
		gtest.Assert(Value.Get(1), "v2")
		gtest.Assert(Value.Get(2), "")
	})
}

