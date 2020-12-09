// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gconv"
	"reflect"
)

// Json converts <r> to JSON format content.
func (r Record) Json() string {
	content, _ := gparser.VarToJson(r.Map())
	return gconv.UnsafeBytesToStr(content)
}

// Xml converts <r> to XML format content.
func (r Record) Xml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.Map(), rootTag...)
	return gconv.UnsafeBytesToStr(content)
}

// Map converts <r> to map[string]interface{}.
func (r Record) Map() Map {
	m := make(map[string]interface{})
	for k, v := range r {
		m[k] = v.Val()
	}
	return m
}

// GMap converts <r> to a gmap.
func (r Record) GMap() *gmap.StrAnyMap {
	return gmap.NewStrAnyMapFrom(r.Map())
}

// Struct converts <r> to a struct.
// Note that the parameter <pointer> should be type of *struct/**struct.
//
// Note that it returns sql.ErrNoRows if <r> is empty.
func (r Record) Struct(pointer interface{}) error {
	// If the record is empty, it returns error.
	if r.IsEmpty() {
		return sql.ErrNoRows
	}
	// Special handling for parameter type: reflect.Value
	if _, ok := pointer.(reflect.Value); ok {
		return convertMapToStruct(r.Map(), pointer)
	}
	var (
		reflectValue = reflect.ValueOf(pointer)
		reflectKind  = reflectValue.Kind()
	)
	if reflectKind != reflect.Ptr {
		return gerror.New("parameter should be type of *struct/**struct")
	}
	reflectValue = reflectValue.Elem()
	reflectKind = reflectValue.Kind()
	if reflectKind == reflect.Invalid {
		return gerror.New("parameter is an invalid pointer, maybe nil")
	}
	if reflectKind != reflect.Ptr && reflectKind != reflect.Struct {
		return gerror.New("parameter should be type of *struct/**struct")
	}
	return convertMapToStruct(r.Map(), pointer)
}

// IsEmpty checks and returns whether <r> is empty.
func (r Record) IsEmpty() bool {
	return len(r) == 0
}
