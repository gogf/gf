// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package gcmd provides console operations, like options/values reading and command running.
package gcmd

import (
	"os"
	"regexp"

	"github.com/gogf/gf/g/os/glog"
)

// Console values.
type gCmdValue struct {
	values []string
}

// Console options.
type gCmdOption struct {
	options map[string]string
}

var Value = &gCmdValue{}                 // Console values.
var Option = &gCmdOption{}               // Console options.
var cmdFuncMap = make(map[string]func()) // Registered callback functions.

func init() {
	doInit()
}

// doInit does the initialization for this package.
func doInit() {
	Value.values = Value.values[:0]
	Option.options = make(map[string]string)
	reg := regexp.MustCompile(`^\-{1,2}(\w+)={0,1}(.*)`)
	for i := 0; i < len(os.Args); i++ {
		result := reg.FindStringSubmatch(os.Args[i])
		if len(result) > 1 {
			Option.options[result[1]] = result[2]
		} else {
			Value.values = append(Value.values, os.Args[i])
		}
	}
}

// BindHandle registers callback function <f> with <cmd>.
func BindHandle(cmd string, f func()) {
	if _, ok := cmdFuncMap[cmd]; ok {
		glog.Fatal("duplicated handle for command:" + cmd)
	} else {
		cmdFuncMap[cmd] = f
	}
}

// RunHandle executes the callback function registered by <cmd>.
func RunHandle(cmd string) {
	if handle, ok := cmdFuncMap[cmd]; ok {
		handle()
	} else {
		glog.Fatal("no handle found for command:" + cmd)
	}
}

// AutoRun automatically recognizes and executes the callback function
// by value of index 0 (the first console parameter).
func AutoRun() {
	if cmd := Value.Get(1); cmd != "" {
		if handle, ok := cmdFuncMap[cmd]; ok {
			handle()
		} else {
			glog.Fatal("no handle found for command:" + cmd)
		}
	} else {
		glog.Fatal("no command found")
	}
}
