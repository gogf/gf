// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gvalid

import "strings"

// 校验错误对象
type Error map[string]map[string]string

// 只获取第一个校验错误项
func (e Error) FirstItem() (string, map[string]string) {
    for k, m := range e {
        return k, m
    }
    return "", nil
}

// 只获取第一个校验错误项的规则及错误信息
func (e Error) FirstRule() (string, string) {
    for _, m := range e {
        for k, v := range m {
            return k, v
        }
    }
    return "", ""
}

// 只获取第一个校验错误项的错误信息
func (e Error) FirstString() (string) {
    for _, m := range e {
        for _, v := range m {
            return v
        }
    }
    return ""
}

// 将所有错误信息构建称字符串，多个错误信息字符串使用"; "符号分隔
func (e Error) String() string {
    return strings.Join(e.Strings(), "; ")
}

// 只返回错误信息，构造成字符串数组返回
func (e Error) Strings() []string {
    array := make([]string, 0)
    for _, m := range e {
        for _, v := range m {
            array = append(array, v)
        }
    }
    return array
}