// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// BindHandle registers callback function `f` with `cmd`.
func (p *Parser) BindHandle(cmd string, f func()) error {
	if _, ok := p.commandFuncMap[cmd]; ok {
		return gerror.NewCode(gcode.CodeInvalidOperation, "duplicated handle for command:"+cmd)
	} else {
		p.commandFuncMap[cmd] = f
	}
	return nil
}

// BindHandleMap registers callback function with map `m`.
func (p *Parser) BindHandleMap(m map[string]func()) error {
	var err error
	for k, v := range m {
		if err = p.BindHandle(k, v); err != nil {
			return err
		}
	}
	return err
}

// RunHandle executes the callback function registered by `cmd`.
func (p *Parser) RunHandle(cmd string) error {
	if handle, ok := p.commandFuncMap[cmd]; ok {
		handle()
	} else {
		return gerror.NewCode(gcode.CodeMissingConfiguration, "no handle found for command:"+cmd)
	}
	return nil
}

// AutoRun automatically recognizes and executes the callback function
// by value of index 0 (the first console parameter).
func (p *Parser) AutoRun() error {
	if cmd := p.GetArg(1); !cmd.IsEmpty() {
		if handle, ok := p.commandFuncMap[cmd.String()]; ok {
			handle()
		} else {
			return gerror.NewCode(gcode.CodeMissingConfiguration, "no handle found for command:"+cmd.String())
		}
	} else {
		return gerror.NewCode(gcode.CodeMissingParameter, "no command found")
	}
	return nil
}
