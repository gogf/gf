// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"

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
//	| _int2, _int4                 | []int32       | // Note: pgtype uses int32 for compatibility
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
//	| PostgreSQL Type | SQL Type                       | pgtype Type              | Go Type     |
//	|-----------------|--------------------------------|--------------------------|-------------|
//	| int2            | int2, smallint                 | -                        | int         |
//	| int4            | int4, integer                  | -                        | int         |
//	| int8            | int8, bigint, bigserial        | -                        | int64       |
//	| uuid            | uuid                           | -                        | uuid.UUID   |
//	| _int2           | int2[], smallint[]             | pgtype.Array[int32]      | []int32     |
//	| _int4           | int4[], integer[]              | pgtype.Array[int32]      | []int32     |
//	| _int8           | int8[], bigint[]               | pgtype.Array[int64]      | []int64     |
//	| _float4         | float4[], real[]               | pgtype.Array[float32]    | []float32   |
//	| _float8         | float8[], double precision[]   | pgtype.Array[float64]    | []float64   |
//	| _bool           | boolean[], bool[]              | pgtype.Array[bool]       | []bool      |
//	| _varchar        | varchar[], character varying[] | pgtype.Array[string]     | []string    |
//	| _text           | text[]                         | pgtype.Array[string]     | []string    |
//	| _char, _bpchar  | char[], character[]            | pgtype.Array[string]     | []string    |
//	| _numeric        | numeric[]                      | pgtype.Array[float64]    | []float64   |
//	| _decimal        | decimal[]                      | pgtype.Array[float64]    | []float64   |
//	| _money          | money[]                        | pgtype.Array[float64]    | []float64   |
//	| _bytea          | bytea[]                        | pgtype.Array[[]byte]     | [][]byte    |
//	| _uuid           | uuid[]                         | pgtype.Array[string]     | []uuid.UUID |
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
		return scanInt32Array(fieldValue)

	// []int64
	case "_int8":
		return scanInt64Array(fieldValue)

	// []float32
	case "_float4":
		return scanFloat32Array(fieldValue)

	// []float64
	case "_float8":
		return scanFloat64Array(fieldValue)

	// []bool
	case "_bool":
		return scanBoolArray(fieldValue)

	// []string
	case "_varchar", "_text", "_char", "_bpchar":
		return scanStringArray(fieldValue)

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
		strArray, err := scanStringArray(fieldValue)
		if err != nil {
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
		return scanFloat64Array(fieldValue)

	// [][]byte
	case "_bytea":
		return scanByteaArray(fieldValue)

	default:
		return d.Core.ConvertValueForLocal(ctx, fieldType, fieldValue)
	}
}

// parsePostgresArray parses PostgreSQL array text format (e.g., "{1,2,3}") into string slice.
// It handles NULL values, quoted strings, and escaped characters.
func parsePostgresArray(src any) ([]string, error) {
	var str string
	switch v := src.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return nil, fmt.Errorf("unsupported type for PostgreSQL array: %T", src)
	}

	if str == "" || str == "{}" {
		return []string{}, nil
	}

	// Validate PostgreSQL array format: must start with '{' and end with '}'
	if !strings.HasPrefix(str, "{") || !strings.HasSuffix(str, "}") {
		return nil, fmt.Errorf("invalid PostgreSQL array format: %s", str)
	}

	// Remove outer braces
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	var result []string
	var current strings.Builder
	var inQuotes bool
	var escaped bool

	for i := 0; i < len(str); i++ {
		c := str[i]
		if escaped {
			current.WriteByte(c)
			escaped = false
			continue
		}
		switch c {
		case '\\':
			escaped = true
		case '"':
			inQuotes = !inQuotes
		case ',':
			if inQuotes {
				current.WriteByte(c)
			} else {
				result = append(result, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	// Add the last element
	if current.Len() > 0 || len(result) > 0 {
		result = append(result, current.String())
	}

	return result, nil
}

// scanInt32Array parses PostgreSQL int2[] or int4[] array.
func scanInt32Array(fieldValue any) ([]int32, error) {
	elements, err := parsePostgresArray(fieldValue)
	if err != nil {
		return nil, err
	}
	result := make([]int32, 0, len(elements))
	for _, s := range elements {
		if s == "NULL" || s == "" {
			continue
		}
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, int32(v))
	}
	return result, nil
}

// scanInt64Array parses PostgreSQL int8[] array.
func scanInt64Array(fieldValue any) ([]int64, error) {
	elements, err := parsePostgresArray(fieldValue)
	if err != nil {
		return nil, err
	}
	result := make([]int64, 0, len(elements))
	for _, s := range elements {
		if s == "NULL" || s == "" {
			continue
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

// scanFloat32Array parses PostgreSQL float4[] array.
func scanFloat32Array(fieldValue any) ([]float32, error) {
	elements, err := parsePostgresArray(fieldValue)
	if err != nil {
		return nil, err
	}
	result := make([]float32, 0, len(elements))
	for _, s := range elements {
		if s == "NULL" || s == "" {
			continue
		}
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, float32(v))
	}
	return result, nil
}

// scanFloat64Array parses PostgreSQL float8[], numeric[], decimal[], money[] array.
func scanFloat64Array(fieldValue any) ([]float64, error) {
	elements, err := parsePostgresArray(fieldValue)
	if err != nil {
		return nil, err
	}
	result := make([]float64, 0, len(elements))
	for _, s := range elements {
		if s == "NULL" || s == "" {
			continue
		}
		// Handle money format like "$1,234.56"
		s = strings.ReplaceAll(s, "$", "")
		s = strings.ReplaceAll(s, ",", "")
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

// scanBoolArray parses PostgreSQL bool[] array.
func scanBoolArray(fieldValue any) ([]bool, error) {
	elements, err := parsePostgresArray(fieldValue)
	if err != nil {
		return nil, err
	}
	result := make([]bool, 0, len(elements))
	for _, s := range elements {
		if s == "NULL" || s == "" {
			continue
		}
		v := s == "t" || s == "true" || s == "TRUE" || s == "1"
		result = append(result, v)
	}
	return result, nil
}

// scanStringArray parses PostgreSQL varchar[], text[], char[] array.
func scanStringArray(fieldValue any) ([]string, error) {
	return parsePostgresArray(fieldValue)
}

// scanByteaArray parses PostgreSQL bytea[] array.
func scanByteaArray(fieldValue any) ([][]byte, error) {
	elements, err := parsePostgresArray(fieldValue)
	if err != nil {
		return nil, err
	}
	result := make([][]byte, 0, len(elements))
	for _, s := range elements {
		if s == "NULL" || s == "" {
			result = append(result, nil)
			continue
		}
		// PostgreSQL bytea is typically in hex format: \x...
		if strings.HasPrefix(s, "\\x") {
			decoded, err := hex.DecodeString(s[2:])
			if err != nil {
				return nil, err
			}
			result = append(result, decoded)
		} else {
			result = append(result, []byte(s))
		}
	}
	return result, nil
}
