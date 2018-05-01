// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import (
    "gitee.com/johng/gf/g/util/gconv"
)

// 将结果列表按照指定的字段值做map[string]Map
func (r Result) ToStringMap(key string) map[string]Record {
    m := make(map[string]Record)
    for _, item := range r {
        if v, ok := item[key]; ok {
            m[gconv.String(v)] = item
        }
    }
    return m
}

// 将结果列表按照指定的字段值做map[int]Map
func (r Result) ToIntMap(key string) map[int]Record {
    m := make(map[int]Record)
    for _, item := range r {
        if v, ok := item[key]; ok {
            m[gconv.Int(v)] = item
        }
    }
    return m
}

// 将结果列表按照指定的字段值做map[uint]Map
func (r Result) ToUintMap(key string) map[uint]Record {
    m := make(map[uint]Record)
    for _, item := range r {
        if v, ok := item[key]; ok {
            m[gconv.Uint(v)] = item
        }
    }
    return m
}
