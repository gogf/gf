// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"context"
	"io"
)

// MustShell performs as Shell, but it panics if any error occurs.
func MustShell(ctx context.Context, cmd string, out io.Writer, in io.Reader) {
	if err := Shell(ctx, cmd, out, in); err != nil {
		panic(err)
	}
}

// MustShellRun performs as ShellRun, but it panics if any error occurs.
func MustShellRun(ctx context.Context, cmd string) {
	if err := ShellRun(ctx, cmd); err != nil {
		panic(err)
	}
}

// MustShellExec performs as ShellExec, but it panics if any error occurs.
func MustShellExec(ctx context.Context, cmd string, environment ...[]string) string {
	result, err := ShellExec(ctx, cmd, environment...)
	if err != nil {
		panic(err)
	}
	return result
}
