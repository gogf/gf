// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

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
