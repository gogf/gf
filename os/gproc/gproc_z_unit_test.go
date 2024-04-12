// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gproc_test

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_ShellExec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s, err := gproc.ShellExec(gctx.New(), `echo 123`)
		t.AssertNil(err)
		t.Assert(s, "123\n")
	})
	// error
	gtest.C(t, func(t *gtest.T) {
		_, err := gproc.ShellExec(gctx.New(), `NoneExistCommandCall`)
		t.AssertNE(err, nil)
	})
}

func Test_ProcessStart(t *testing.T) {
	var ctx = gctx.GetInitCtx()
	gtest.C(t, func(t *gtest.T) {
		var buf strings.Builder
		// Necessary check.
		protoc := gproc.SearchBinary("gf")
		t.AssertNE(protoc, "")
		var command = gproc.NewProcess(protoc, nil)
		command.Args = append(command.Args, "version")
		command.Stdout = &buf
		err := command.Run(ctx)
		t.AssertNil(err)
		//	v2.7.0
		//	Welcome to GoFrame!
		//		Env Detail:
		t.Assert(strings.Contains(buf.String(), "Welcome to GoFrame!\nEnv Detail:"), true)
	})

}
