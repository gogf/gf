// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"

	"github.com/gogf/gf/test/gtest"
)

func Test_BuildOptions(t *testing.T) {
	gtest.Case(t, func() {
		s := gcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		})
		gtest.Assert(s, "-n=john")
	})

	gtest.Case(t, func() {
		s := gcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		}, "-test")
		gtest.Assert(s, "-testn=john")
	})
}
