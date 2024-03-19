// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

//go:build windows

package gproc

import (
	"syscall"

	"github.com/gogf/gf/v2/text/gstr"
)

// Because when the underlying parameters are passed in on the Windows platform,
// escape characters will be added, causing some commands to fail.
func newProcess(p *Process, args []string, path string) *Process {
	p.SysProcAttr = &syscall.SysProcAttr{}
	p.SysProcAttr.CmdLine = path + " " + gstr.Join(args, " ")
	return p
}
