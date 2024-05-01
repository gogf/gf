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
	"strconv"
)

func RegisterFieldConverterFunc(ctx context.Context, db DB,
	tableField *sql.ColumnType, structField reflect.StructField) (convertFn fieldConvertFunc, tempArg any) {
	tempArg = &sql.RawBytes{}

	tableFieldIsTimeType := func(typ reflect.Type) bool {
		switch typ.String() {
		case "time.Time":
			return true
		case "sql.NullTime":
			return true
		default:
			// todo typ.ConvertibleTo(reflect.TypeOf(time.Time{}))
			return false
		}
	}
	tabFieldType := tableField
	localType, _ := db.CheckLocalTypeForField(ctx, tableField.DatabaseTypeName(), nil)

	switch localType {
	case LocalTypeBytes:
		convertFn = getStringConvertFunc(structField.Type)

	case LocalTypeString:
		// To support string to other types of conversions, it must be compatible
		// Cannot be converted to a numeric type
		convertFn = getStringConvertFunc(structField.Type)
	case LocalTypeInt, LocalTypeInt64:
		convertFn = getIntegerConvertFunc[int64](structField.Type, strconv.ParseInt)
	case LocalTypeUint, LocalTypeUint64:
		convertFn = getIntegerConvertFunc[uint64](structField.Type, strconv.ParseUint)
	case LocalTypeFloat32, LocalTypeFloat64:
		convertFn = getFloatConvertFunc(structField.Type)
	case LocalTypeBool:
		convertFn = getBoolConvertFunc(structField.Type)
	case LocalTypeDate:
		// The time format of different databases may be different, and the compatibility is not good enough
		// Unified use of sql.RawBytes are received and processed by standard libraries or drivers
		ok := tableFieldIsTimeType(tabFieldType.ScanType())
		if ok {
			tempArg = &sql.NullTime{}
		}
		convertFn = getTimeConvertFunc(structField.Type)
	case LocalTypeDatetime:
		// The time format of different databases may be different, and the compatibility is not good enough
		// Unified use of sql.RawBytes are received and processed by standard libraries or drivers
		ok := tableFieldIsTimeType(tabFieldType.ScanType())
		if ok {
			tempArg = &sql.NullTime{}
		}
		convertFn = getTimeConvertFunc(structField.Type)
	case LocalTypeDecimal: // float
		convertFn = getDecimalConvertFunc(structField.Type)
	case LocalTypeInt64Bytes:
		convertFn = getBitConvertFunc(structField.Type)
	case LocalTypeJson:
		convertFn = getJsonConvertFunc(structField.Type, true)
	case LocalTypeJsonb:
		convertFn = getJsonConvertFunc(structField.Type, true)
	default:
	}
	if convertFn == nil {
		panic(&typeConvertError{
			driverName:  db.GetConfig().Type,
			columnName:  tabFieldType.Name(),
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
