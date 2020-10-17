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
	gtest.C(t, func(t *gtest.T) {
		os.Args = []string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}
		t.Assert(len(gcmd.GetArgAll()), 4)
		t.Assert(gcmd.GetArg(1), "remove")
		t.Assert(gcmd.GetArg(100, "test"), "test")
		t.Assert(gcmd.GetOpt("n"), "")
		t.Assert(gcmd.ContainsOpt("p"), true)
		t.Assert(gcmd.ContainsOpt("n"), true)
		t.Assert(gcmd.ContainsOpt("none"), false)
		t.Assert(gcmd.GetOpt("none", "value"), "value")

	})
}

func Test_BuildOptions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		})
		t.Assert(s, "-n=john")
	})

	gtest.C(t, func(t *gtest.T) {
		s := gcmd.BuildOptions(g.MapStrStr{
			"n": "john",
		}, "-test")
		t.Assert(s, "-testn=john")
	})
}
