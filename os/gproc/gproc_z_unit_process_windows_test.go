// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

//go:build windows

package gproc_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_ProcessRun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gf := gproc.SearchBinary("gf")
		if gf == "" {
			return
		}
		var command = gproc.NewProcess(gf, nil)
		command.Args = append(command.Args, "version")
		var buf strings.Builder
		command.Stdout = &buf
		command.Stderr = &buf
		err := command.Run(gctx.GetInitCtx())
		t.AssertNil(err)

		errOutput := `up         upgrade GoFrame version/tool to latest one in current project`
		t.Assert(gstr.Contains(buf.String(), errOutput), false)
	})

	gtest.C(t, func(t *gtest.T) {
		binary := gproc.SearchBinary("go")
		t.AssertNE(binary, "")
		var command = gproc.NewProcess(binary, nil)
		command.Args = append(command.Args, "version")
		var buf strings.Builder
		command.Stdout = &buf
		command.Stderr = &buf
		err := command.Run(gctx.GetInitCtx())
		t.AssertNil(err)

		errOutput := `bug         start a bug report`
		t.Assert(gstr.Contains(buf.String(), errOutput), false)
	})

	gtest.C(t, func(t *gtest.T) {
		binary := gproc.SearchBinary("go")
		t.AssertNE(binary, "")
		var command = gproc.NewProcess(binary, nil)

		testpath := gtest.DataPath("gobuild")
		filename := filepath.Join(testpath, "main.go")
		output := filepath.Join(testpath, "main.exe")

		command.Args = append(command.Args, "build")
		command.Args = append(command.Args, `-ldflags="-X 'main.TestString=\"test string\"'"`)
		command.Args = append(command.Args, "-o", output)
		command.Args = append(command.Args, filename)

		var buf strings.Builder
		command.Stdout = &buf
		command.Stderr = &buf
		err := command.Run(gctx.GetInitCtx())
		t.AssertNil(err)

		realPath, err := gfile.Search(output)
		t.AssertNE(realPath, "")
		t.Assert(err, nil)
		defer gfile.Remove(output)

		result, err := gproc.ShellExec(gctx.New(), output)
		t.Assert(err, nil)
		t.Assert(gstr.Contains(result, "test string"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		binary := gproc.SearchBinary("go")
		t.AssertNE(binary, "")
		var command = gproc.NewProcess(binary, nil)

		testpath := gtest.DataPath("gobuild")
		filename := filepath.Join(testpath, "main.go")
		output := filepath.Join(testpath, "main.exe")

		command.Args = append(command.Args, "build")
		command.Args = append(command.Args, `-ldflags="-s -w"`)
		command.Args = append(command.Args, "-o", output)
		command.Args = append(command.Args, filename)

		err := command.Run(gctx.GetInitCtx())
		t.AssertNil(err)

		realPath, err := gfile.Search(output)
		t.Assert(err, nil)
		t.AssertNE(realPath, "")

		defer gfile.Remove(output)
	})
}
