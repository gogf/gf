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

	"github.com/gogf/gf/container/garray"

	"github.com/gogf/gf/os/gcmd"

	"github.com/gogf/gf/test/gtest"
)

func Test_Parse(t *testing.T) {
	gtest.Case(t, func() {
		os.Args = []string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}
		p, err := gcmd.Parse(map[string]bool{
			"n, name":   true,
			"p, prefix": true,
			"f,force":   false,
			"q,quiet":   false,
		})
		gtest.Assert(err, nil)
		gtest.Assert(len(p.GetArgAll()), 3)
		gtest.Assert(p.GetArg(0), "gf")
		gtest.Assert(p.GetArg(1), "remove")
		gtest.Assert(p.GetArg(2), "path")
		gtest.Assert(p.GetArgVar(2).String(), "path")

		gtest.Assert(len(p.GetOptAll()), 8)
		gtest.Assert(p.GetOpt("n"), "root")
		gtest.Assert(p.GetOpt("name"), "root")
		gtest.Assert(p.GetOpt("p"), "www")
		gtest.Assert(p.GetOpt("prefix"), "www")
		gtest.Assert(p.GetOptVar("prefix").String(), "www")

		gtest.Assert(p.ContainsOpt("n"), true)
		gtest.Assert(p.ContainsOpt("name"), true)
		gtest.Assert(p.ContainsOpt("p"), true)
		gtest.Assert(p.ContainsOpt("prefix"), true)
		gtest.Assert(p.ContainsOpt("f"), true)
		gtest.Assert(p.ContainsOpt("force"), true)
		gtest.Assert(p.ContainsOpt("q"), true)
		gtest.Assert(p.ContainsOpt("quiet"), true)
		gtest.Assert(p.ContainsOpt("none"), false)
	})
}

func Test_ParseWithArgs(t *testing.T) {
	gtest.Case(t, func() {
		p, err := gcmd.ParseWithArgs(
			[]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"},
			map[string]bool{
				"n, name":   true,
				"p, prefix": true,
				"f,force":   false,
				"q,quiet":   false,
			})
		gtest.Assert(err, nil)
		gtest.Assert(len(p.GetArgAll()), 3)
		gtest.Assert(p.GetArg(0), "gf")
		gtest.Assert(p.GetArg(1), "remove")
		gtest.Assert(p.GetArg(2), "path")
		gtest.Assert(p.GetArgVar(2).String(), "path")

		gtest.Assert(len(p.GetOptAll()), 8)
		gtest.Assert(p.GetOpt("n"), "root")
		gtest.Assert(p.GetOpt("name"), "root")
		gtest.Assert(p.GetOpt("p"), "www")
		gtest.Assert(p.GetOpt("prefix"), "www")
		gtest.Assert(p.GetOptVar("prefix").String(), "www")

		gtest.Assert(p.ContainsOpt("n"), true)
		gtest.Assert(p.ContainsOpt("name"), true)
		gtest.Assert(p.ContainsOpt("p"), true)
		gtest.Assert(p.ContainsOpt("prefix"), true)
		gtest.Assert(p.ContainsOpt("f"), true)
		gtest.Assert(p.ContainsOpt("force"), true)
		gtest.Assert(p.ContainsOpt("q"), true)
		gtest.Assert(p.ContainsOpt("quiet"), true)
		gtest.Assert(p.ContainsOpt("none"), false)
	})
}

func Test_Handler(t *testing.T) {
	gtest.Case(t, func() {
		p, err := gcmd.ParseWithArgs(
			[]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"},
			map[string]bool{
				"n, name":   true,
				"p, prefix": true,
				"f,force":   false,
				"q,quiet":   false,
			})
		gtest.Assert(err, nil)
		array := garray.New()
		err = p.BindHandle("remove", func() {
			array.Append(1)
		})
		gtest.Assert(err, nil)

		err = p.BindHandle("remove", func() {
			array.Append(1)
		})
		gtest.AssertNE(err, nil)

		err = p.BindHandle("test", func() {
			array.Append(1)
		})
		gtest.Assert(err, nil)

		err = p.RunHandle("remove")
		gtest.Assert(err, nil)
		gtest.Assert(array.Len(), 1)

		err = p.RunHandle("none")
		gtest.AssertNE(err, nil)
		gtest.Assert(array.Len(), 1)

		err = p.RunHandle("test")
		gtest.Assert(err, nil)
		gtest.Assert(array.Len(), 2)

		err = p.AutoRun()
		gtest.Assert(err, nil)
		gtest.Assert(array.Len(), 3)
	})
}
