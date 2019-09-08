// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"

	"github.com/gogf/gf/test/gtest"
)

func Test_Default(t *testing.T) {
	gtest.Case(t, func() {
		os.Args = []string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}
		gtest.Assert(len(gcmd.GetArgAll()), 4)
		gtest.Assert(gcmd.GetArg(1), "remove")
		gtest.Assert(gcmd.GetArg(100, "test"), "test")
		gtest.Assert(gcmd.GetOpt("n"), "")
		gtest.Assert(gcmd.ContainsOpt("p"), true)
		gtest.Assert(gcmd.ContainsOpt("n"), true)
		gtest.Assert(gcmd.ContainsOpt("none"), false)
		gtest.Assert(gcmd.GetOpt("none", "value"), "value")

	})
}

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
