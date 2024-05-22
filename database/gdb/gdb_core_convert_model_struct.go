// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"reflect"
)

func (t *Table) makeStructQueryModel(columns int) *queryStructModel {
	scaArgs := make([]any, columns)
	q := &queryStructModel{
		table:    t,
		scanArgs: scaArgs,
	}
	for i := 0; i < columns; i++ {
		// The reason for doing this is because the [queryStructModel] implements the [sql.Scanner] interface
		// This way, when calling the standard library's [sql.Rows.Scan], you can enter our custom conversion logic
		// Control which field needs to be assigned the current value to through the [queryStructModel.scanIndex] variable
		q.scanArgs[i] = q
	}
	return q
}

// Leave it for subsequent extensions,
// if you pass the struct pointer directly without taking the map
type queryStructModel struct {
	table     *Table
	scanArgs  []any
	scanIndex int
	Struct    reflect.Value
}

func (q *queryStructModel) Scan(src any) error {
	field := q.table.fields[q.scanIndex]
	if field.convertFunc == nil {
		// Indicates that this field is redundant and does not exist in the struct
		q.scanIndex++
		return nil
	}
	fieldValue := field.GetReflectValue(q.Struct)
	err := field.convertFunc(fieldValue, src)
	if err != nil {
		err = fmt.Errorf("it is not possible to convert from `%v :%T`(%s: %s) to `%s: %s` err:%v",
			src, src,
			field.ColumnFieldName, field.ColumnFieldType.DatabaseTypeName(),
			field.StructField.Name, field.StructFieldType, err)
	}
	q.scanIndex++
	return err
}

func (q *queryStructModel) next(structValue reflect.Value) {
	q.Struct = structValue
	// The index needs to be initialized
	q.scanIndex = 0
}
