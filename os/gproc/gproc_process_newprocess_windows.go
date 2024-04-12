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

// The Windows platform goes back directly and does nothing
// When the Process.Start method is called, it is handled in joinProcessArgs
func newProcess(p *Process, _ []string, _ string) *Process {
	return p
}

// When the Process.Start method is called,
// it will be called on the Windows platform
func joinProcessArgs(p *Process) {
	p.SysProcAttr = &syscall.SysProcAttr{}
	p.SysProcAttr.CmdLine = gstr.Join(p.Args, " ")
}
