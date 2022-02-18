// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"io"
)

// MustShell performs as Shell, but it panics if any error occurs.
func MustShell(cmd string, out io.Writer, in io.Reader) {
	if err := Shell(cmd, out, in); err != nil {
		panic(err)
	}
}

// MustShellRun performs as ShellRun, but it panics if any error occurs.
func MustShellRun(cmd string) {
	if err := ShellRun(cmd); err != nil {
		panic(err)
	}
}

// MustShellExec performs as ShellExec, but it panics if any error occurs.
func MustShellExec(cmd string, environment ...[]string) string {
	result, err := ShellExec(cmd, environment...)
	if err != nil {
		panic(err)
	}
	return result
}
