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
		// 表示这个字段是多余的，在结构体中不存在
		q.scanIndex++
		return nil
	}

	fieldValue := reflect.New(field.StructFieldType).Elem()
	err := field.convertFunc(fieldValue, src)
	if err != nil {
		err = fmt.Errorf("不能从`%v:%T`(%s: %s)转换到`%s: %s` err: %v",
			src, src,
			field.ColumnField, field.ColumnFieldType.DatabaseTypeName(),
			field.StructField.Name, field.StructFieldType, err)
	}
	q.scanIndex++
	q.Map[field.ColumnField] = gvar.New(fieldValue.Interface())
	return err
}

func (q *queryMapModel) next(mapValue map[string]Value) {
	q.Map = mapValue
	// 索引需要初始化
	q.scanIndex = 0
}
