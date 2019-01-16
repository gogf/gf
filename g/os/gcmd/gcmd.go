// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

// Package gcmd provides console operations, like options/values reading and command running.
// 
// 命令行管理.
package gcmd

import (
	"os"
    "errors"
    "regexp"
)

// 命令行参数列表
type gCmdValue  struct {
    values []string
}

// 命令行选项列表
type gCmdOption struct {
    options map[string]string
}

// 终端管理对象(全局)
var Value      = &gCmdValue{}            // 终端参数-命令参数列表
var Option     = &gCmdOption{}           // 终端参数-选项参数列表
var cmdFuncMap = make(map[string]func()) // 终端命令及函数地址对应表

// 检查并初始化console参数，在包加载的时候触发
// 初始化时执行，不影响运行时性能
func init() {
    reg := regexp.MustCompile(`\-\-{0,1}(.+?)=(.+)`)
    Option.options = make(map[string]string)
    for i := 0; i < len(os.Args); i++ {
        result := reg.FindStringSubmatch(os.Args[i])
        if len(result) > 1 {
            Option.options[result[1]] = result[2]
        } else {
            Value.values = append(Value.values, os.Args[i])
        }
    }
}

// 返回所有的命令行参数values
func (c *gCmdValue) GetAll() []string {
    return c.values
}

// 返回所有的命令行参数options
func (c *gCmdOption) GetAll() map[string]string {
    return c.options
}

// 获得一条指定索引位置的value参数
func (c *gCmdValue) Get(index uint8, def...string) string {
    if index < uint8(len(c.values)) {
        return c.values[index]
    } else if len(def) > 0 {
        return def[0]
    }
    return ""
}

// 获得一条指定索引位置的option参数;
func (c *gCmdOption) Get(key string, def...string) string {
    if option, ok := c.options[key]; ok {
        return option
    } else if len(def) > 0 {
        return def[0]
    }
    return ""
}

// 绑定命令行参数及对应的命令函数，注意命令函数参数是函数的内存地址
// 如果操作失败返回错误信息
func BindHandle (cmd string, f func()) error {
    if _, ok := cmdFuncMap[cmd]; ok {
        return errors.New("duplicated handle for command:" + cmd)
    } else {
        cmdFuncMap[cmd] = f
        return nil
    }
}

// 执行命令对应的函数
func RunHandle (cmd string) error {
    if handle, ok := cmdFuncMap[cmd]; ok {
        handle()
        return nil
    } else {
        return errors.New("no handle found for command:" + cmd)
    }
}

// 自动识别命令参数并执行命令参数对应的函数
func AutoRun () error {
    if cmd := Value.Get(1); cmd != "" {
        if handle, ok := cmdFuncMap[cmd]; ok {
            handle()
            return nil
        } else {
            return errors.New("no handle found for command:" + cmd)
        }
    } else {
        return errors.New("no command found")
    }
}
