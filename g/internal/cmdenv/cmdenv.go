// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmdenv

import (
	"github.com/gogf/gf/g/container/gvar"
	"os"
	"regexp"
	"strings"
)

var (
	// Console options.
	cmdOptions = make(map[string]string)
)

func init() {
	reg := regexp.MustCompile(`\-\-{0,1}(.+?)=(.+)`)
	for i := 0; i < len(os.Args); i++ {
		result := reg.FindStringSubmatch(os.Args[i])
		if len(result) > 1 {
			cmdOptions[result[1]] = result[2]
		}
	}
}

// 获取指定名称的命令行参数，当不存在时获取环境变量参数，皆不存在时，返回给定的默认值。
// 规则:
// 1、命令行参数以小写字母格式，使用: gf.包名.变量名 传递；
// 2、环境变量参数以大写字母格式，使用: GF_包名_变量名 传递；
func Get(key string, def...interface{}) *gvar.Var {
    value := interface{}(nil)
    if len(def) > 0 {
        value = def[0]
    }
    if v, ok := cmdOptions[key]; ok {
        value = v
    } else {
        key = strings.ToUpper(strings.Replace(key, ".", "_", -1))
        if v := os.Getenv(key); v != "" {
            value = v
        }
    }
    return gvar.New(value, true)
}
