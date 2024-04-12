// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

//go:build windows

package gproc_test

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
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
}
