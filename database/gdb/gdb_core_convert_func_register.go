// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

func registerFieldConvertFunc(ctx context.Context, db DB,
	tableField *sql.ColumnType, structField reflect.StructField) (convertFn fieldConvertFunc) {
	localType, _ := db.CheckLocalTypeForField(ctx, tableField.DatabaseTypeName(), nil)
	// 有几个特殊的类型，需要特殊处理
	switch localType {
	case LocalTypeUint64Bytes:
		// bit
		convertFn = getBitConvertFunc(structField.Type, 0)
	case LocalTypeInt64Bytes:
		// bit
		convertFn = getBitConvertFunc(structField.Type, 0)
	case LocalTypeDecimal:
		// decimal numeric money
		convertFn = getDecimalConvertFunc(structField.Type, 0)
	default:
		convertFn = getConverter(structField.Type, 0)
	}
	if convertFn == nil {
		panic(&typeConvertError{
			driverName:  db.GetConfig().Type,
			columnName:  tableField.Name(),
			columnType:  tableField.DatabaseTypeName(),
			structField: structField,
		})
	}
	return
}

type typeConvertError struct {
	driverName  string
	columnName  string
	columnType  string
	structField reflect.StructField
}

func (t *typeConvertError) Error() string {
	err := `Driver: %s does not support conversion from (%s: %s) to (%s: %s)`
	return fmt.Sprintf(err, t.driverName, t.columnName, t.columnType, t.structField.Name, t.structField.Type)
}
