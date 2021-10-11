// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcmd_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/v2/container/garray"

	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Parse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		os.Args = []string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}
		p, err := gcmd.Parse(map[string]bool{
			"n, name":   true,
			"p, prefix": true,
			"f,force":   false,
			"q,quiet":   false,
		})
		t.Assert(err, nil)
		t.Assert(len(p.GetArgAll()), 3)
		t.Assert(p.GetArg(0), "gf")
		t.Assert(p.GetArg(1), "remove")
		t.Assert(p.GetArg(2), "path")
		t.Assert(p.GetArgVar(2).String(), "path")

		t.Assert(len(p.GetOptAll()), 8)
		t.Assert(p.GetOpt("n"), "root")
		t.Assert(p.GetOpt("name"), "root")
		t.Assert(p.GetOpt("p"), "www")
		t.Assert(p.GetOpt("prefix"), "www")
		t.Assert(p.GetOptVar("prefix").String(), "www")

		t.Assert(p.ContainsOpt("n"), true)
		t.Assert(p.ContainsOpt("name"), true)
		t.Assert(p.ContainsOpt("p"), true)
		t.Assert(p.ContainsOpt("prefix"), true)
		t.Assert(p.ContainsOpt("f"), true)
		t.Assert(p.ContainsOpt("force"), true)
		t.Assert(p.ContainsOpt("q"), true)
		t.Assert(p.ContainsOpt("quiet"), true)
		t.Assert(p.ContainsOpt("none"), false)
	})
}

func Test_ParseWithArgs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p, err := gcmd.ParseWithArgs(
			[]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"},
			map[string]bool{
				"n, name":   true,
				"p, prefix": true,
				"f,force":   false,
				"q,quiet":   false,
			})
		t.Assert(err, nil)
		t.Assert(len(p.GetArgAll()), 3)
		t.Assert(p.GetArg(0), "gf")
		t.Assert(p.GetArg(1), "remove")
		t.Assert(p.GetArg(2), "path")
		t.Assert(p.GetArgVar(2).String(), "path")

		t.Assert(len(p.GetOptAll()), 8)
		t.Assert(p.GetOpt("n"), "root")
		t.Assert(p.GetOpt("name"), "root")
		t.Assert(p.GetOpt("p"), "www")
		t.Assert(p.GetOpt("prefix"), "www")
		t.Assert(p.GetOptVar("prefix").String(), "www")

		t.Assert(p.ContainsOpt("n"), true)
		t.Assert(p.ContainsOpt("name"), true)
		t.Assert(p.ContainsOpt("p"), true)
		t.Assert(p.ContainsOpt("prefix"), true)
		t.Assert(p.ContainsOpt("f"), true)
		t.Assert(p.ContainsOpt("force"), true)
		t.Assert(p.ContainsOpt("q"), true)
		t.Assert(p.ContainsOpt("quiet"), true)
		t.Assert(p.ContainsOpt("none"), false)
	})
}

func Test_Handler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p, err := gcmd.ParseWithArgs(
			[]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"},
			map[string]bool{
				"n, name":   true,
				"p, prefix": true,
				"f,force":   false,
				"q,quiet":   false,
			})
		t.Assert(err, nil)
		array := garray.New()
		err = p.BindHandle("remove", func() {
			array.Append(1)
		})
		t.Assert(err, nil)

		err = p.BindHandle("remove", func() {
			array.Append(1)
		})
		t.AssertNE(err, nil)

		err = p.BindHandle("test", func() {
			array.Append(1)
		})
		t.Assert(err, nil)

		err = p.RunHandle("remove")
		t.Assert(err, nil)
		t.Assert(array.Len(), 1)

		err = p.RunHandle("none")
		t.AssertNE(err, nil)
		t.Assert(array.Len(), 1)

		err = p.RunHandle("test")
		t.Assert(err, nil)
		t.Assert(array.Len(), 2)

		err = p.AutoRun()
		t.Assert(err, nil)
		t.Assert(array.Len(), 3)
	})
}
