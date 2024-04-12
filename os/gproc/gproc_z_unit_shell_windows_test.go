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
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_ShellExec_GoBuild_Windows(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		testpath := gtest.DataPath("gobuild")
		filename := filepath.Join(testpath, "main.go")
		output := filepath.Join(testpath, "main.exe")
		cmd := fmt.Sprintf(`go build -ldflags="-s -w" -o %s  %s`, output, filename)

		err := gproc.ShellRun(gctx.New(), cmd)
		t.Assert(err, nil)

		realPath, err := gfile.Search(output)
		t.AssertNE(realPath, "")
		t.Assert(err, nil)

		defer gfile.Remove(output)
	})

	gtest.C(t, func(t *gtest.T) {
		testpath := gtest.DataPath("gobuild")
		filename := filepath.Join(testpath, "main.go")
		output := filepath.Join(testpath, "main.exe")
		cmd := fmt.Sprintf(`go build -ldflags="-X 'main.TestString=\"test string\"'" -o %s %s`, output, filename)

		err := gproc.ShellRun(gctx.New(), cmd)
		t.Assert(err, nil)

		realPath, err := gfile.Search(output)
		t.AssertNE(realPath, "")
		t.Assert(err, nil)
		defer gfile.Remove(output)

		result, err := gproc.ShellExec(gctx.New(), output)
		t.Assert(err, nil)
		t.Assert(gstr.Contains(result, "test string"), true)
	})
}
