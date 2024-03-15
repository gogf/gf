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
