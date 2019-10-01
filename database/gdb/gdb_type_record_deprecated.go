// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"github.com/gogf/gf/encoding/gparser"
)

// Deprecated.
func (r Record) ToJson() string {
	content, _ := gparser.VarToJson(r.Map())
	return string(content)
}

// Deprecated.
func (r Record) ToXml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.Map(), rootTag...)
	return string(content)
}

// Deprecated.
func (r Record) ToMap() Map {
	m := make(map[string]interface{})
	for k, v := range r {
		m[k] = v.Val()
	}
	return m
}

// Deprecated.
func (r Record) ToStruct(pointer interface{}) error {
	if r == nil {
		return sql.ErrNoRows
	}
	return mapToStruct(r.Map(), pointer)
}
