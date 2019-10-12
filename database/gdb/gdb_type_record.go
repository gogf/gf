// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"github.com/gogf/gf/container/gmap"

	"github.com/gogf/gf/encoding/gparser"
)

// 将记录结果转换为JSON字符串
func (r Record) Json() string {
	content, _ := gparser.VarToJson(r.Map())
	return string(content)
}

// 将记录结果转换为XML字符串
func (r Record) Xml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.Map(), rootTag...)
	return string(content)
}

// 将Record转换为Map类型
func (r Record) Map() Map {
	m := make(map[string]interface{})
	for k, v := range r {
		m[k] = v.Val()
	}
	return m
}

// 将Record转换为常用的gmap.StrAnyMap类型
func (r Record) GMap() *gmap.StrAnyMap {
	return gmap.NewStrAnyMapFrom(r.Map())
}

// 将Map变量映射到指定的struct对象中，注意参数应当是一个对象的指针
func (r Record) Struct(pointer interface{}) error {
	if r == nil {
		return sql.ErrNoRows
	}
	return mapToStruct(r.Map(), pointer)
}
