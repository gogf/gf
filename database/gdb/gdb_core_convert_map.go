// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gvar"
)

func (c *Core) scanRowsToMap(rows *sql.Rows, table *Table, columns int) (result Result, err error) {
	mapQuery := table.makeMapQueryModel(columns)
	for {
		record := Record{}
		mapQuery.next(record)
		if err = rows.Scan(mapQuery.scanArgs...); err != nil {
			return result, err
		}
		result = append(result, record)
		if !rows.Next() {
			break
		}
	}
	return
}

func (t *Table) makeMapQueryModel(columns int) *queryMapModel {
	scaArgs := make([]any, columns)
	q := &queryMapModel{
		table:    t,
		scanArgs: scaArgs,
	}
	for i := 0; i < columns; i++ {
		// The reason for doing this is because the [queryMapModel] implements the [sql.Scanner] interface
		// This way, when calling the standard library's [sql.Rows.Scan], you can enter our custom conversion logic
		// Control which field needs to be assigned the current value to through the [queryMapModel.scanIndex] variable
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
	if field.convertFunc == nil || src == nil {
		// Indicates that this field is redundant and does not exist in the struct
		q.scanIndex++
		return nil
	}

	fieldValue := reflect.New(field.StructFieldType).Elem()
	err := field.convertFunc(fieldValue, src)
	if err != nil {
		err = fmt.Errorf("it is not possible to convert from `%v :%T`(%s: %s) to `%s: %s` err: %v",
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
