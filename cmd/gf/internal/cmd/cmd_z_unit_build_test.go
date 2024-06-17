// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Build_Single(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			buildPath  = gtest.DataPath(`build`, `single`)
			pwd        = gfile.Pwd()
			binaryName = `t.test`
			binaryPath = gtest.DataPath(`build`, `single`, binaryName)
			f          = cBuild{}
		)
		defer gfile.Chdir(pwd)
		defer gfile.Remove(binaryPath)
		err := gfile.Chdir(buildPath)
		t.AssertNil(err)

		t.Assert(gfile.Exists(binaryPath), false)
		_, err = f.Index(ctx, cBuildInput{
			File: cBuildDefaultFile,
			Name: binaryName,
		})
		t.AssertNil(err)
		t.Assert(gfile.Exists(binaryPath), true)
	})
}

func Test_Build_Single_Output(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			buildPath     = gtest.DataPath(`build`, `single`)
			pwd           = gfile.Pwd()
			binaryName    = `tt`
			binaryDirPath = gtest.DataPath(`build`, `single`, `tt`)
			binaryPath    = gtest.DataPath(`build`, `single`, `tt`, binaryName)
			f             = cBuild{}
		)
		defer gfile.Chdir(pwd)
		defer gfile.Remove(binaryDirPath)
		err := gfile.Chdir(buildPath)
		t.AssertNil(err)

		t.Assert(gfile.Exists(binaryPath), false)
		_, err = f.Index(ctx, cBuildInput{
			Output: "./tt/tt",
		})
		t.AssertNil(err)
		t.Assert(gfile.Exists(binaryPath), true)
	})
}

func Test_Build_Single_Path(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			buildPath     = gtest.DataPath(`build`, `single`)
			pwd           = gfile.Pwd()
			dirName       = "ttt"
			binaryName    = `main`
			binaryDirPath = gtest.DataPath(`build`, `single`, dirName)
			binaryPath    = gtest.DataPath(`build`, `single`, dirName, binaryName)
			f             = cBuild{}
		)
		defer gfile.Chdir(pwd)
		defer gfile.Remove(binaryDirPath)
		err := gfile.Chdir(buildPath)
		t.AssertNil(err)

		t.Assert(gfile.Exists(binaryPath), false)
		_, err = f.Index(ctx, cBuildInput{
			Path: "ttt",
		})
		t.AssertNil(err)
		t.Assert(gfile.Exists(binaryPath), true)
	})
}

func Test_Build_Single_VarMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			buildPath  = gtest.DataPath(`build`, `varmap`)
			pwd        = gfile.Pwd()
			binaryName = `main`
			binaryPath = gtest.DataPath(`build`, `varmap`, binaryName)
			f          = cBuild{}
		)
		defer gfile.Chdir(pwd)
		defer gfile.Remove(binaryPath)
		err := gfile.Chdir(buildPath)
		t.AssertNil(err)

		t.Assert(gfile.Exists(binaryPath), false)
		_, err = f.Index(ctx, cBuildInput{
			VarMap: map[string]interface{}{
				"a": "1",
				"b": "2",
			},
		})
		t.AssertNil(err)
		t.Assert(gfile.Exists(binaryPath), true)

		result, err := gproc.ShellExec(ctx, binaryPath)
		t.AssertNil(err)
		t.Assert(gstr.Contains(result, `a: 1`), true)
		t.Assert(gstr.Contains(result, `b: 2`), true)
	})
}

func Test_Build_Multiple(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			buildPath         = gtest.DataPath(`build`, `multiple`)
			pwd               = gfile.Pwd()
			binaryDirPath     = gtest.DataPath(`build`, `multiple`, `temp`)
			binaryPathLinux   = gtest.DataPath(`build`, `multiple`, `temp`, `v1.1`, `linux_amd64`, `ttt`)
			binaryPathWindows = gtest.DataPath(`build`, `multiple`, `temp`, `v1.1`, `windows_amd64`, `ttt.exe`)
			f                 = cBuild{}
		)
		defer gfile.Chdir(pwd)
		defer gfile.Remove(binaryDirPath)
		err := gfile.Chdir(buildPath)
		t.AssertNil(err)

		t.Assert(gfile.Exists(binaryPathLinux), false)
		t.Assert(gfile.Exists(binaryPathWindows), false)
		_, err = f.Index(ctx, cBuildInput{
			File:    "multiple.go",
			Name:    "ttt",
			Version: "v1.1",
			Arch:    "amd64",
			System:  "linux, windows",
			Path:    "temp",
		})
		t.AssertNil(err)
		t.Assert(gfile.Exists(binaryPathLinux), true)
		t.Assert(gfile.Exists(binaryPathWindows), true)
	})
}
