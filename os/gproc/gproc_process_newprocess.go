//go:build !windows

package gproc

import "strings"

func newProcess(p *Process, args []string, path string) *Process {

	if len(args) > 0 {
		// Exclude of current binary path.
		start := 0
		if strings.EqualFold(path, args[0]) {
			start = 1
		}
		p.Args = append(p.Args, args[start:]...)
	}
	return p
}
