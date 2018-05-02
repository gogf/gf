// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import (
    "gitee.com/johng/gf/g/util/gutil"
)

// 将Record转换为Map，其中最主要的区别是里面的键值被强制转换为string类型，方便json处理
func (r Record) ToMap() Map {
    m := make(map[string]interface{})
    for k, v := range r {
        m[k] = v.String()
    }
    return m
}

// 将Map变量映射到指定的struct对象中，注意参数应当是一个对象的指针
func (r Record) ToStruct(obj interface{}) error {
    m := make(map[string]interface{})
    for k, v := range r {
        m[k] = v
    }
    return gutil.MapToStruct(m, obj)
}
