// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"github.com/gogf/gf/g/encoding/gparser"
)

// 将记录结果转换为JSON字符串
func (r Record) ToJson() string {
    content, _ := gparser.VarToJson(r.ToMap())
    return string(content)
}

// 将记录结果转换为XML字符串
func (r Record) ToXml(rootTag...string) string {
    content, _ := gparser.VarToXml(r.ToMap(), rootTag...)
    return string(content)
}

// 将Record转换为Map，其中最主要的区别是里面的键值被强制转换为string类型，方便json处理
func (r Record) ToMap() Map {
    m := make(map[string]interface{})
    for k, v := range r {
        m[k] = v.Val()
    }
    return m
}

// 将Map变量映射到指定的struct对象中，注意参数应当是一个对象的指针
func (r Record) ToStruct(pointer interface{}) error {
    return mapToStruct(r.ToMap(), pointer)
}
