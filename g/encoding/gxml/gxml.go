// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// XML
package gxml

import (
    "github.com/clbanning/mxj"
)

// 将XML内容解析为map变量
func Decode(xmlbyte []byte) (map[string]interface{}, error) {
    return mxj.NewMapXml(xmlbyte)
}

// 将map变量解析为XML格式内容
func Encode(v map[string]interface{}) ([]byte, error) {
    return mxj.Map(v).Xml()
}

// XML格式内容直接转换为JSON格式内容
func ToJson(xmlbyte []byte) ([]byte, error) {
    if mv, err := mxj.NewMapXml(xmlbyte); err == nil {
        return mv.Json()
    } else {
        return nil, err
    }
}