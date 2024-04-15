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
)

func Test_ProcessRun(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		binary := gproc.SearchBinary("go")
		t.AssertNE(binary, "")
		var command = gproc.NewProcess(binary, nil)

		testPath := gtest.DataPath("gobuild")
		filename := filepath.Join(testPath, "main.go")
		output := filepath.Join(testPath, "main.exe")

		command.Args = append(command.Args, "build")
		command.Args = append(command.Args, `-ldflags="-X 'main.TestString=\"test string\"'"`)
		command.Args = append(command.Args, "-o", output)
		command.Args = append(command.Args, filename)

		err := command.Run(gctx.GetInitCtx())
		t.AssertNil(err)

		exists := gfile.Exists(output)
		t.Assert(exists, true)
		defer gfile.Remove(output)

		runCmd := gproc.NewProcess(output, nil)
		var buf strings.Builder
		runCmd.Stdout = &buf
		runCmd.Stderr = &buf
		err = runCmd.Run(gctx.GetInitCtx())
		t.Assert(err, nil)
		t.Assert(buf.String(), `"test string"`)
	})

	gtest.C(t, func(t *gtest.T) {
		binary := gproc.SearchBinary("go")
		t.AssertNE(binary, "")
		// NewProcess(path,args) path： It's best not to have spaces
		var command = gproc.NewProcess(binary, nil)

		testPath := gtest.DataPath("gobuild")
		filename := filepath.Join(testPath, "main.go")
		output := filepath.Join(testPath, "main.exe")

		command.Args = append(command.Args, "build")
		command.Args = append(command.Args, `-ldflags="-s -w"`)
		command.Args = append(command.Args, "-o", output)
		command.Args = append(command.Args, filename)

		err := command.Run(gctx.GetInitCtx())
		t.AssertNil(err)

		exists := gfile.Exists(output)
		t.Assert(exists, true)

		defer gfile.Remove(output)
	})
}
