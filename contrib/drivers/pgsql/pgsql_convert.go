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

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// ConvertValueForField converts value to database acceptable value.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	if g.IsNil(fieldValue) {
		return d.Core.ConvertValueForField(ctx, fieldType, fieldValue)
	}

	var fieldValueKind = reflect.TypeOf(fieldValue).Kind()

	if fieldValueKind == reflect.Slice {
		// For pgsql, json or jsonb require '[]'
		if !gstr.Contains(fieldType, "json") {
			fieldValue = gstr.ReplaceByMap(gconv.String(fieldValue),
				map[string]string{
					"[": "{",
					"]": "}",
				},
			)
		}
	}
	return d.Core.ConvertValueForField(ctx, fieldType, fieldValue)
}

// CheckLocalTypeForField checks and returns corresponding local golang type for given db type.
// The parameter `fieldType` is in lower case, like:
// `int2`, `int4`, `int8`, `_int2`, `_int4`, `_int8`, `_float4`, `_float8`, etc.
//
// PostgreSQL type mapping:
//
//	| PostgreSQL Type              | Local Go Type |
//	|------------------------------|---------------|
//	| int2, int4                   | int           |
//	| int8                         | int64         |
//	| uuid                         | uuid.UUID     |
//	| _int2, _int4                 | []int32       | // Note: pq package does not provide Int16Array; int32 is used for compatibility
//	| _int8                        | []int64       |
//	| _float4                      | []float32     |
//	| _float8                      | []float64     |
//	| _bool                        | []bool        |
//	| _varchar, _text              | []string      |
//	| _char, _bpchar               | []string      |
//	| _numeric, _decimal, _money   | []float64     |
//	| _bytea                       | [][]byte      |
//	| _uuid                        | []uuid.UUID   |
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (gdb.LocalType, error) {
	var typeName string
	match, _ := gregex.MatchString(`(.+?)\((.+)\)`, fieldType)
	if len(match) == 3 {
		typeName = gstr.Trim(match[1])
	} else {
		typeName = fieldType
	}
	typeName = strings.ToLower(typeName)
	switch typeName {
	case "int2", "int4":
		return gdb.LocalTypeInt, nil

	case "int8":
		return gdb.LocalTypeInt64, nil

	case "uuid":
		return gdb.LocalTypeUUID, nil

	case "_int2", "_int4":
		return gdb.LocalTypeInt32Slice, nil

	case "_int8":
		return gdb.LocalTypeInt64Slice, nil

	case "_float4":
		return gdb.LocalTypeFloat32Slice, nil

	case "_float8":
		return gdb.LocalTypeFloat64Slice, nil

	case "_bool":
		return gdb.LocalTypeBoolSlice, nil

	case "_varchar", "_text", "_char", "_bpchar":
		return gdb.LocalTypeStringSlice, nil

	case "_uuid":
		return gdb.LocalTypeUUIDSlice, nil

	case "_numeric", "_decimal", "_money":
		return gdb.LocalTypeFloat64Slice, nil

	case "_bytea":
		return gdb.LocalTypeBytesSlice, nil

	default:
		return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	}
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
// The parameter `fieldType` is in lower case, like:
// `int2`, `int4`, `int8`, `_int2`, `_int4`, `_int8`, `uuid`, `_uuid`, etc.
//
// See: https://www.postgresql.org/docs/current/datatype.html
//
// PostgreSQL type mapping:
//
//	| PostgreSQL Type | SQL Type                       | pq Type         | Go Type     |
//	|-----------------|--------------------------------|-----------------|-------------|
//	| int2            | int2, smallint                 | -               | int         |
//	| int4            | int4, integer                  | -               | int         |
//	| int8            | int8, bigint, bigserial        | -               | int64       |
//	| uuid            | uuid                           | -               | uuid.UUID   |
//	| _int2           | int2[], smallint[]             | pq.Int32Array   | []int32     |
//	| _int4           | int4[], integer[]              | pq.Int32Array   | []int32     |
//	| _int8           | int8[], bigint[]               | pq.Int64Array   | []int64     |
//	| _float4         | float4[], real[]               | pq.Float32Array | []float32   |
//	| _float8         | float8[], double precision[]   | pq.Float64Array | []float64   |
//	| _bool           | boolean[], bool[]              | pq.BoolArray    | []bool      |
//	| _varchar        | varchar[], character varying[] | pq.StringArray  | []string    |
//	| _text           | text[]                         | pq.StringArray  | []string    |
//	| _char, _bpchar  | char[], character[]            | pq.StringArray  | []string    |
//	| _numeric        | numeric[]                      | pq.Float64Array | []float64   |
//	| _decimal        | decimal[]                      | pq.Float64Array | []float64   |
//	| _money          | money[]                        | pq.Float64Array | []float64   |
//	| _bytea          | bytea[]                        | pq.ByteaArray   | [][]byte    |
//	| _uuid           | uuid[]                         | pq.StringArray  | []uuid.UUID |
//
// Note: PostgreSQL also supports these array types but they are not yet mapped:
//   - _date (date[]), _timestamp (timestamp[]), _timestamptz (timestamptz[])
//   - _jsonb (jsonb[]), _json (json[])
func (d *Driver) ConvertValueForLocal(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	typeName, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	typeName = strings.ToLower(typeName)

	// Basic types are mostly handled by Core layer, only handle array types here
	switch typeName {

	// []int32
	case "_int2", "_int4":
		var result pq.Int32Array
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return []int32(result), nil

	// []int64
	case "_int8":
		var result pq.Int64Array
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return []int64(result), nil

	// []float32
	case "_float4":
		var result pq.Float32Array
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return []float32(result), nil

	// []float64
	case "_float8":
		var result pq.Float64Array
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return []float64(result), nil

	// []bool
	case "_bool":
		var result pq.BoolArray
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return []bool(result), nil

	// []string
	case "_varchar", "_text", "_char", "_bpchar":
		var result pq.StringArray
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return []string(result), nil

	// uuid.UUID
	case "uuid":
		var uuidStr string
		switch v := fieldValue.(type) {
		case []byte:
			uuidStr = string(v)
		case string:
			uuidStr = v
		default:
			uuidStr = gconv.String(fieldValue)
		}
		result, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, err
		}
		return result, nil

	// []uuid.UUID
	case "_uuid":
		var strArray pq.StringArray
		if err := strArray.Scan(fieldValue); err != nil {
			return nil, err
		}
		result := make([]uuid.UUID, len(strArray))
		for i, s := range strArray {
			parsed, err := uuid.Parse(s)
			if err != nil {
				return nil, err
			}
			result[i] = parsed
		}
		return result, nil

	// []float64
	case "_numeric", "_decimal", "_money":
		var result pq.Float64Array
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return []float64(result), nil

	// [][]byte
	case "_bytea":
		var result pq.ByteaArray
		if err := result.Scan(fieldValue); err != nil {
			return nil, err
		}
		return [][]byte(result), nil

	default:
		return d.Core.ConvertValueForLocal(ctx, fieldType, fieldValue)
	}
}
