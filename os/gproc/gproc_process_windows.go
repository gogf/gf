//go:build windows

package gproc

import (
	"syscall"

	"github.com/gogf/gf/v2/text/gstr"
)

func newProcess(p *Process, args []string, path string) *Process {
	p.SysProcAttr = &syscall.SysProcAttr{}
	p.SysProcAttr.CmdLine = path + " " + gstr.Join(args, " ")

	return p
}
