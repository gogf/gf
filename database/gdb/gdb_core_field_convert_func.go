// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// todo  需要返回一个函数
func (c *Core) getFieldConvertFunc(ctx context.Context, value interface{}, columnType *sql.ColumnType) (interface{}, error) {
	var scanType = columnType.ScanType()
	if scanType != nil {
		// Common basic builtin types.
		switch scanType.Kind() {
		case
			reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return gconv.Convert(
				gconv.String(value),
				columnType.ScanType().String(),
			), nil
		}
	}
	// Other complex types, especially custom types.
	return c.getFieldValueConvertFunc(ctx, columnType.DatabaseTypeName(), value)
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
// The parameter `fieldType` is in lower case, like:
// `float(5,2)`, `unsigned double(5,2)`, `decimal(10,2)`, `char(45)`, `varchar(100)`, etc.
func (c *Core) getFieldValueConvertFunc(
	ctx context.Context, fieldType string, fieldValue interface{},
) (interface{}, error) {
	// If there's no type retrieved, it returns the `fieldValue` directly
	// to use its original data type, as `fieldValue` is type of interface{}.
	if fieldType == "" {
		return fieldValue, nil
	}
	typeName, err := c.db.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	if err != nil {
		return nil, err
	}
	switch typeName {
	case LocalTypeBytes:
		var typeNameStr = string(typeName)
		if strings.Contains(typeNameStr, "binary") || strings.Contains(typeNameStr, "blob") {
			return fieldValue, nil
		}
		return gconv.Bytes(fieldValue), nil

	case LocalTypeInt:
		return gconv.Int(gconv.String(fieldValue)), nil

	case LocalTypeUint:
		return gconv.Uint(gconv.String(fieldValue)), nil

	case LocalTypeInt64:
		return gconv.Int64(gconv.String(fieldValue)), nil

	case LocalTypeUint64:
		return gconv.Uint64(gconv.String(fieldValue)), nil

	case LocalTypeInt64Bytes:
		return gbinary.BeDecodeToInt64(gconv.Bytes(fieldValue)), nil

	case LocalTypeUint64Bytes:
		return gbinary.BeDecodeToUint64(gconv.Bytes(fieldValue)), nil

	case LocalTypeFloat32:
		return gconv.Float32(gconv.String(fieldValue)), nil

	case LocalTypeFloat64:
		return gconv.Float64(gconv.String(fieldValue)), nil

	case LocalTypeBool:
		s := gconv.String(fieldValue)
		// mssql is true|false string.
		if strings.EqualFold(s, "true") {
			return 1, nil
		}
		if strings.EqualFold(s, "false") {
			return 0, nil
		}
		return gconv.Bool(fieldValue), nil

	case LocalTypeDate:
		// Date without time.
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t).Format("Y-m-d"), nil
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t.Format("Y-m-d"), nil

	case LocalTypeDatetime:
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t), nil
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t, nil

	default:
		return gconv.String(fieldValue), nil
	}
}
