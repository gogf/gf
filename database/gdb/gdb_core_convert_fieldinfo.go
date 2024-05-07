// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"reflect"
)

type fieldConvertFunc = func(dst reflect.Value, src any) error

type fieldConvertInfo struct {
	// table field
	ColumnField      string
	ColumnFieldIndex int
	ColumnFieldType  *sql.ColumnType
	// struct field
	StructField      reflect.StructField
	StructFieldType  reflect.Type
	StructFieldIndex []int
	convertFunc      fieldConvertFunc
}

// GetReflectValue = reflect.Value.FieldByIndex
func (c *fieldConvertInfo) GetReflectValue(structValue reflect.Value) reflect.Value {
	if len(c.StructFieldIndex) == 1 {
		return structValue.Field(c.StructFieldIndex[0])
	}
	v := structValue
	for i, x := range c.StructFieldIndex {
		if i > 0 {
			if v.Kind() == reflect.Pointer {
				if v.IsNil() {
					v.Set(reflect.New(v.Type().Elem()))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}
