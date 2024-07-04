// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"context"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// ConvertValueForField converts value to database acceptable value.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) {
	var (
		fieldValueKind = reflect.TypeOf(fieldValue).Kind()
	)

	if fieldValueKind == reflect.Slice {
		fieldValue = gstr.ReplaceByMap(gconv.String(fieldValue),
			map[string]string{
				"[": "{",
				"]": "}",
			},
		)
	}
	return d.Core.ConvertValueForField(ctx, fieldType, fieldValue)
}

// CheckLocalTypeForField checks and returns corresponding local golang type for given db type.
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue interface{}) (gdb.LocalType, error) {
	var typeName string
	match, _ := gregex.MatchString(`(.+?)\((.+)\)`, fieldType)
	if len(match) == 3 {
		typeName = gstr.Trim(match[1])
	} else {
		typeName = fieldType
	}
	typeName = strings.ToLower(typeName)
	switch typeName {
	case
		// For pgsql, int2 = smallint.
		"int2",
		// For pgsql, int4 = integer
		"int4":
		return gdb.LocalTypeInt, nil

	case
		// For pgsql, int8 = bigint
		"int8":
		return gdb.LocalTypeInt64, nil

	case
		"_int2",
		"_int4":
		return gdb.LocalTypeIntSlice, nil

	case
		"_int8":
		return gdb.LocalTypeInt64Slice, nil

	default:
		return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	}
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
// The parameter `fieldType` is in lower case, like:
// `float(5,2)`, `unsigned double(5,2)`, `decimal(10,2)`, `char(45)`, `varchar(100)`, etc.
func (d *Driver) ConvertValueForLocal(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) {
	typeName, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	typeName = strings.ToLower(typeName)
	switch typeName {
	// For pgsql, int2 = smallint and int4 = integer.
	case "int2", "int4":
		return gconv.Int(gconv.String(fieldValue)), nil

	// For pgsql, int8 = bigint.
	case "int8":
		return gconv.Int64(gconv.String(fieldValue)), nil

	// Int32 slice.
	case
		"_int2", "_int4":
		return gconv.Ints(
			gstr.ReplaceByMap(gconv.String(fieldValue),
				map[string]string{
					"{": "[",
					"}": "]",
				},
			),
		), nil

	// Int64 slice.
	case
		"_int8":
		return gconv.Int64s(
			gstr.ReplaceByMap(gconv.String(fieldValue),
				map[string]string{
					"{": "[",
					"}": "]",
				},
			),
		), nil

	default:
		return d.Core.ConvertValueForLocal(ctx, fieldType, fieldValue)
	}
}
