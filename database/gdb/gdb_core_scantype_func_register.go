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

// 默认实现
func RegisterFieldConverterFunc(ctx context.Context, db DB,
	tableField *sql.ColumnType, structField reflect.StructField) (convertFn fieldScanFunc, tempArg any) {
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

	localType, err := db.CheckLocalTypeForField(ctx, tableField.DatabaseTypeName(), nil)
	_ = err

	switch localType {
	case LocalTypeBytes:
		convertFn = getStringConvertFunc(structField.Type)

	case LocalTypeString:
		// todo 支持string到其他类型的转换，必须是兼容的
		// 不能转换到数字类型的
		convertFn = getStringConvertFunc(structField.Type)

	case LocalTypeInt:

		convertFn = getIntegerConvertFunc[int64](structField.Type, strconv.ParseInt)

	case LocalTypeUint:
		convertFn = getIntegerConvertFunc[uint64](structField.Type, strconv.ParseUint)

	case LocalTypeInt64:
		convertFn = getIntegerConvertFunc[int64](structField.Type, strconv.ParseInt)

	case LocalTypeUint64:
		convertFn = getIntegerConvertFunc[uint64](structField.Type, strconv.ParseUint)

	case LocalTypeFloat32:
		convertFn = getFloatConvertFunc(structField.Type)

	case LocalTypeFloat64:
		convertFn = getFloatConvertFunc(structField.Type)

	case LocalTypeBool:
		convertFn = getBoolConvertFunc(structField.Type)

	case LocalTypeDate:
		// 可能不同数据库存储的时间格式不一样，后续如果兼容性不够好的话
		// 统一使用sql.RawBytes接收，让标准库或者驱动去处理
		ok := tableFieldIsTimeType(tabFieldType.ScanType())
		if ok {
			tempArg = &sql.NullTime{}
		}

		convertFn = getTimeConvertFunc(structField.Type)

	case LocalTypeDatetime:
		// 可能不同数据库存储的时间格式不一样，后续如果兼容性不够好的话
		// 统一使用sql.RawBytes接收，让标准库或者驱动去处理
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
		panic(fmt.Errorf("不支持的类型  字段名:%s 表字段类型:%v localType:%v",
			tableField.Name(), tableField.DatabaseTypeName(),
			localType,
		))
	}
	return
}
