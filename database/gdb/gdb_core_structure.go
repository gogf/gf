// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// ConvertDataForRecord is a very important function, which does converting for any data that
// will be inserted into table/collection as a record.
//
// The parameter `value` should be type of *map/map/*struct/struct.
// It supports embedded struct definition for struct.
func (c *Core) ConvertDataForRecord(ctx context.Context, value interface{}) map[string]interface{} {
	var data = DataToMapDeep(value)
	for k, v := range data {
		data[k] = c.ConvertDataForRecordValue(ctx, v)
	}
	return data
}

func (c *Core) ConvertDataForRecordValue(ctx context.Context, value interface{}) interface{} {
	var (
		rvValue = reflect.ValueOf(value)
		rvKind  = rvValue.Kind()
	)
	for rvKind == reflect.Ptr {
		rvValue = rvValue.Elem()
		rvKind = rvValue.Kind()
	}
	switch rvKind {
	case reflect.Slice, reflect.Array, reflect.Map:
		// It should ignore the bytes type.
		if _, ok := value.([]byte); !ok {
			// Convert the value to JSON.
			value, _ = json.Marshal(value)
		}

	case reflect.Struct:
		switch r := value.(type) {
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		case time.Time:
			if r.IsZero() {
				value = nil
			}

		case gtime.Time:
			if r.IsZero() {
				value = nil
			}

		case *gtime.Time:
			if r.IsZero() {
				value = nil
			}

		case *time.Time:
			// Nothing to do.

		case Counter, *Counter:
			// Nothing to do.

		default:
			// Use string conversion in default.
			if s, ok := value.(iString); ok {
				value = s.String()
			} else {
				// Convert the value to JSON.
				value, _ = json.Marshal(value)
			}
		}
	}
	return value
}

// convertFieldValueToLocalValue automatically checks and converts field value from database type
// to golang variable type as underlying value of Value.
func (c *Core) convertFieldValueToLocalValue(fieldValue interface{}, fieldType string) interface{} {
	// If there's no type retrieved, it returns the `fieldValue` directly
	// to use its original data type, as `fieldValue` is type of interface{}.
	if fieldType == "" {
		return fieldValue
	}
	typeName, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	typeName = strings.ToLower(typeName)
	switch typeName {
	case
		"binary",
		"varbinary",
		"blob",
		"tinyblob",
		"mediumblob",
		"longblob":
		return gconv.Bytes(fieldValue)

	case
		"int",
		"tinyint",
		"small_int",
		"smallint",
		"medium_int",
		"mediumint",
		"serial":
		if gstr.ContainsI(fieldType, "unsigned") {
			gconv.Uint(gconv.String(fieldValue))
		}
		return gconv.Int(gconv.String(fieldValue))

	case
		"int8", // For pgsql, int8 = bigint.
		"big_int",
		"bigint",
		"bigserial":
		if gstr.ContainsI(fieldType, "unsigned") {
			gconv.Uint64(gconv.String(fieldValue))
		}
		return gconv.Int64(gconv.String(fieldValue))

	case "real":
		return gconv.Float32(gconv.String(fieldValue))

	case
		"float",
		"double",
		"decimal",
		"money",
		"numeric",
		"smallmoney":
		return gconv.Float64(gconv.String(fieldValue))

	case "bit":
		s := gconv.String(fieldValue)
		// mssql is true|false string.
		if strings.EqualFold(s, "true") {
			return 1
		}
		if strings.EqualFold(s, "false") {
			return 0
		}
		return gbinary.BeDecodeToInt64(gconv.Bytes(fieldValue))

	case "bool":
		return gconv.Bool(fieldValue)

	case "date":
		// Date without time.
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t).Format("Y-m-d")
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t.Format("Y-m-d")

	case
		"datetime",
		"timestamp",
		"timestamptz":
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t)
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t

	default:
		// Auto-detect field type, using key match.
		switch {
		case strings.Contains(typeName, "text") || strings.Contains(typeName, "char") || strings.Contains(typeName, "character"):
			return gconv.String(fieldValue)

		case strings.Contains(typeName, "float") || strings.Contains(typeName, "double") || strings.Contains(typeName, "numeric"):
			return gconv.Float64(gconv.String(fieldValue))

		case strings.Contains(typeName, "bool"):
			return gconv.Bool(gconv.String(fieldValue))

		case strings.Contains(typeName, "binary") || strings.Contains(typeName, "blob"):
			return fieldValue

		case strings.Contains(typeName, "int"):
			return gconv.Int(gconv.String(fieldValue))

		case strings.Contains(typeName, "time"):
			s := gconv.String(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t

		case strings.Contains(typeName, "date"):
			s := gconv.String(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t

		default:
			return gconv.String(fieldValue)
		}
	}
}

// mappingAndFilterData automatically mappings the map key to table field and removes
// all key-value pairs that are not the field of given table.
func (c *Core) mappingAndFilterData(schema, table string, data map[string]interface{}, filter bool) (map[string]interface{}, error) {
	fieldsMap, err := c.TableFields(c.guessPrimaryTableName(table), schema)
	if err != nil {
		return nil, err
	}
	fieldsKeyMap := make(map[string]interface{}, len(fieldsMap))
	for k := range fieldsMap {
		fieldsKeyMap[k] = nil
	}
	// Automatic data key to table field name mapping.
	var foundKey string
	for dataKey, dataValue := range data {
		if _, ok := fieldsKeyMap[dataKey]; !ok {
			foundKey, _ = gutil.MapPossibleItemByKey(fieldsKeyMap, dataKey)
			if foundKey != "" {
				if _, ok = data[foundKey]; !ok {
					data[foundKey] = dataValue
				}
				delete(data, dataKey)
			}
		}
	}
	// Data filtering.
	// It deletes all key-value pairs that has incorrect field name.
	if filter {
		for dataKey := range data {
			if _, ok := fieldsMap[dataKey]; !ok {
				delete(data, dataKey)
			}
		}
	}
	return data, nil
}
