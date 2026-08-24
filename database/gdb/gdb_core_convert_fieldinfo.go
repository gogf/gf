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

// dest represents the structural field to be assigned a value to
// src is the field value of the database
type fieldConvertFunc = func(dest reflect.Value, src any) error

type fieldConvertInfo struct {
	// table field
	ColumnFieldName  string
	ColumnFieldIndex int
	ColumnFieldType  *sql.ColumnType
	// struct field
	StructField     reflect.StructField
	StructFieldType reflect.Type
	// The reason why an index is an []int is because this field may be an anonymous structure
	StructFieldIndex []int
	convertFunc      fieldConvertFunc
}

// GetReflectValue = reflect.Value.FieldByIndex
func (c *fieldConvertInfo) GetReflectValue(structValue reflect.Value) reflect.Value {
	if len(c.StructFieldIndex) == 1 {
		return structValue.Field(c.StructFieldIndex[0])
	}

	fieldValue := structValue.Field(c.StructFieldIndex[0])
	for i := 1; i < len(c.StructFieldIndex); i++ {
		if fieldValue.Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}
			fieldValue = fieldValue.Elem()
		}
		fieldValue = fieldValue.Field(c.StructFieldIndex[i])
	}
	return fieldValue
}
