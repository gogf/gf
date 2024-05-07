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
		q.scanArgs[i] = q
	}
	return q
}

// 留作后续的扩展，如果直接传结构体指针不走map的话
type queryStructModel struct {
	table     *Table
	scanArgs  []any
	scanIndex int
	Struct    reflect.Value
}

func (q *queryStructModel) Scan(src any) error {
	fieldName := q.table.fieldsIndex[q.scanIndex]
	field := q.table.fieldsMap[fieldName]
	if field.convertFunc == nil {
		// 表示这个字段是多余的，在结构体中不存在
		q.scanIndex++
		return nil
	}
	fieldValue := field.GetReflectValue(q.Struct)
	err := field.convertFunc(fieldValue, src)
	q.scanIndex++
	if err != nil {
		err = fmt.Errorf("不能从`%v: %T`(%s: %s)转换到`%s: %s` err: %v",
			src, src,
			field.ColumnField, field.ColumnFieldType.DatabaseTypeName(),
			field.StructField.Name, field.StructFieldType, err)
	}
	return err
}

func (q *queryStructModel) next(structValue reflect.Value) {
	q.Struct = structValue
	// 索引需要初始化
	q.scanIndex = 0
}
