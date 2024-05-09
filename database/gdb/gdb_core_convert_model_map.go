// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
)

func (t *Table) makeMapQueryModel(columns int) *queryMapModel {
	scaArgs := make([]any, columns)
	q := &queryMapModel{
		table:    t,
		scanArgs: scaArgs,
	}
	for i := 0; i < columns; i++ {
		q.scanArgs[i] = q
	}
	return q
}

type queryMapModel struct {
	table     *Table
	scanArgs  []any
	scanIndex int
	Map       map[string]Value
}

func (q *queryMapModel) Scan(src any) error {
	field := q.table.fields[q.scanIndex]
	if field.convertFunc == nil {
		// Indicates that this field is redundant and does not exist in the struct
		q.scanIndex++
		return nil
	}

	fieldValue := reflect.New(field.StructFieldType).Elem()
	err := field.convertFunc(fieldValue, src)
	if err != nil {
		err = fmt.Errorf("it is not possible to convert from `%v :%T`(%s: %s) to `%s: %s` err:%v",
			src, src,
			field.ColumnFieldName, field.ColumnFieldType.DatabaseTypeName(),
			field.StructField.Name, field.StructFieldType, err)
	}
	q.scanIndex++
	q.Map[field.ColumnFieldName] = gvar.New(fieldValue.Interface())
	return err
}

func (q *queryMapModel) next(mapValue map[string]Value) {
	q.Map = mapValue
	// The index needs to be initialized
	q.scanIndex = 0
}
