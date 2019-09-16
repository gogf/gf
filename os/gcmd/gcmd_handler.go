// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"errors"
)

// BindHandle registers callback function <f> with <cmd>.
func BindHandle(cmd string, f func()) error {
	if _, ok := defaultCommandFuncMap[cmd]; ok {
		return errors.New("duplicated handle for command:" + cmd)
	} else {
		defaultCommandFuncMap[cmd] = f
	}
	return nil
}

// BindHandle registers callback function with map <m>.
func BindHandleMap(m map[string]func()) error {
	var err error
	for k, v := range m {
		if err = BindHandle(k, v); err != nil {
			return err
		}
	}
	return err
}

// RunHandle executes the callback function registered by <cmd>.
func RunHandle(cmd string) error {
	if handle, ok := defaultCommandFuncMap[cmd]; ok {
		handle()
	} else {
		return errors.New("no handle found for command:" + cmd)
	}
	return nil
}

// AutoRun automatically recognizes and executes the callback function
// by value of index 0 (the first console parameter).
func AutoRun() error {
	if cmd := GetArg(1); cmd != "" {
		if handle, ok := defaultCommandFuncMap[cmd]; ok {
			handle()
		} else {
			return errors.New("no handle found for command:" + cmd)
		}
	} else {
		return errors.New("no command found")
	}
	return nil
}
