// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package cmdenv

import (
    "gitee.com/johng/gf/g/container/gvar"
    "gitee.com/johng/gf/g/os/gcmd"
    "gitee.com/johng/gf/g/os/genv"
    "strings"
)

// 获取指定名称的命令行参数，当不存在时获取环境变量参数，皆不存在时，返回给定的默认值。
// 规则:
// 1、命令行参数以小写字母格式，使用: gf.包名.变量名 传递；
// 2、环境变量参数以大写字母格式，使用: GF_包名_变量名 传递；
func Get(key string, def...interface{}) *gvar.Var {
    value := interface{}(nil)
    if len(def) > 0 {
        value = def[0]
    }
    if v := gcmd.Option.Get(key); v != "" {
        value = v
    } else {
        key = strings.ToUpper(strings.Replace(key, ".", "_", -1))
        if v := genv.Get(key); v != "" {
            value = v
        }
    }
    return gvar.New(value, true)
}
