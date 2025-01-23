// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

//go:build windows

package gproc_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_ShellExec_GoBuild_Windows(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		testPath := gtest.DataPath("gobuild")
		filename := filepath.Join(testPath, "main.go")
		output := filepath.Join(testPath, "main.exe")
		cmd := fmt.Sprintf(`go build -ldflags="-s -w" -o %s  %s`, output, filename)

		err := gproc.ShellRun(gctx.New(), cmd)
		t.Assert(err, nil)

		exists := gfile.Exists(output)
		t.Assert(exists, true)

		defer gfile.Remove(output)
	})

	gtest.C(t, func(t *gtest.T) {
		testPath := gtest.DataPath("gobuild")
		filename := filepath.Join(testPath, "main.go")
		output := filepath.Join(testPath, "main.exe")
		cmd := fmt.Sprintf(`go build -ldflags="-X 'main.TestString=\"test string\"'" -o %s %s`, output, filename)

		err := gproc.ShellRun(gctx.New(), cmd)
		t.Assert(err, nil)

		exists := gfile.Exists(output)
		t.Assert(exists, true)
		defer gfile.Remove(output)

		result, err := gproc.ShellExec(gctx.New(), output)
		t.Assert(err, nil)
		t.Assert(result, `"test string"`)
	})

}

func Test_ShellExec_SpaceDir_Windows(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		testPath := gtest.DataPath("shellexec")
		filename := filepath.Join(testPath, "main.go")
		// go build -o test.exe main.go
		cmd := fmt.Sprintf(`go build -o test.exe %s`, filename)
		r, err := gproc.ShellExec(gctx.New(), cmd)
		t.AssertNil(err)
		t.Assert(r, "")

		exists := gfile.Exists(filename)
		t.Assert(exists, true)

		outputDir := filepath.Join(testPath, "testdir")
		output := filepath.Join(outputDir, "test.exe")
		err = gfile.Move("test.exe", output)
		t.AssertNil(err)
		defer gfile.Remove(output)

		expectContent := "123"
		testOutput := filepath.Join(testPath, "space dir", "test.txt")
		cmd = fmt.Sprintf(`%s -c %s -o "%s"`, output, expectContent, testOutput)
		r, err = gproc.ShellExec(gctx.New(), cmd)
		t.AssertNil(err)

		exists = gfile.Exists(testOutput)
		t.Assert(exists, true)
		defer gfile.Remove(testOutput)

		contents := gfile.GetContents(testOutput)
		t.Assert(contents, expectContent)
	})
	gtest.C(t, func(t *gtest.T) {
		testPath := gtest.DataPath("shellexec")
		filename := filepath.Join(testPath, "main.go")
		// go build -o test.exe main.go
		cmd := fmt.Sprintf(`go build -o test.exe %s`, filename)
		r, err := gproc.ShellExec(gctx.New(), cmd)
		t.AssertNil(err)
		t.Assert(r, "")

		exists := gfile.Exists(filename)
		t.Assert(exists, true)

		outputDir := filepath.Join(testPath, "space dir")
		output := filepath.Join(outputDir, "test.exe")
		err = gfile.Move("test.exe", output)
		t.AssertNil(err)
		defer gfile.Remove(output)

		expectContent := "123"
		testOutput := filepath.Join(testPath, "testdir", "test.txt")
		cmd = fmt.Sprintf(`"%s" -c %s -o %s`, output, expectContent, testOutput)
		r, err = gproc.ShellExec(gctx.New(), cmd)
		t.AssertNil(err)

		exists = gfile.Exists(testOutput)
		t.Assert(exists, true)
		defer gfile.Remove(testOutput)

		contents := gfile.GetContents(testOutput)
		t.Assert(contents, expectContent)

	})
}
