// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"github.com/gogf/gf/os/genv"
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"

	"github.com/gogf/gf/test/gtest"
)

func Test_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gcmd.Init([]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}...)
		t.Assert(len(gcmd.GetArgAll()), 2)
		t.Assert(gcmd.GetArg(1), "path")
		t.Assert(gcmd.GetArg(100, "test"), "test")
		t.Assert(gcmd.GetOpt("force"), "remove")
		t.Assert(gcmd.GetOpt("n"), "root")
		t.Assert(gcmd.ContainsOpt("fq"), true)
		t.Assert(gcmd.ContainsOpt("p"), true)
		t.Assert(gcmd.ContainsOpt("none"), false)
		t.Assert(gcmd.GetOpt("none", "value"), "value")
	})
	gtest.C(t, func(t *gtest.T) {
		gcmd.Init([]string{"gf", "gen", "-h"}...)
		t.Assert(len(gcmd.GetArgAll()), 2)
		t.Assert(gcmd.GetOpt("h"), "")
		t.Assert(gcmd.ContainsOpt("h"), true)
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

func Test_GetWithEnv(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		genv.Set("TEST", "1")
		defer genv.Remove("TEST")
		t.Assert(gcmd.GetWithEnv("test"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		genv.Set("TEST", "1")
		defer genv.Remove("TEST")
		gcmd.Init("-test", "2")
		t.Assert(gcmd.GetWithEnv("test"), 2)
	})
}
