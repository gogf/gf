// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gutil"
	"strings"
	"time"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/os/gtime"

	"github.com/gogf/gf/encoding/gbinary"

	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// convertFieldValueToLocalValue automatically checks and converts field value from database type
// to golang variable type.
func (c *Core) convertFieldValueToLocalValue(fieldValue interface{}, fieldType string) interface{} {
	// If there's no type retrieved, it returns the <fieldValue> directly
	// to use its original data type, as <fieldValue> is type of interface{}.
	if fieldType == "" {
		return fieldValue
	}
	t, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	t = strings.ToLower(t)
	switch t {
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
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t).Format("Y-m-d")
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t.Format("Y-m-d")

	case
		"datetime",
		"timestamp":
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t)
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t.String()

	default:
		// Auto detect field type, using key match.
		switch {
		case strings.Contains(t, "text") || strings.Contains(t, "char") || strings.Contains(t, "character"):
			return gconv.String(fieldValue)

		case strings.Contains(t, "float") || strings.Contains(t, "double") || strings.Contains(t, "numeric"):
			return gconv.Float64(gconv.String(fieldValue))

		case strings.Contains(t, "bool"):
			return gconv.Bool(gconv.String(fieldValue))

		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			return fieldValue

		case strings.Contains(t, "int"):
			return gconv.Int(gconv.String(fieldValue))

		case strings.Contains(t, "time"):
			s := gconv.String(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t.String()

		case strings.Contains(t, "date"):
			s := gconv.String(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t.Format("Y-m-d")

		default:
			return gconv.String(fieldValue)
		}
	}
}

// mappingAndFilterData automatically mappings the map key to table field and removes
// all key-value pairs that are not the field of given table.
func (c *Core) mappingAndFilterData(schema, table string, data map[string]interface{}, filter bool) (map[string]interface{}, error) {
	if fieldsMap, err := c.DB.TableFields(table, schema); err == nil {
		fieldsKeyMap := make(map[string]interface{}, len(fieldsMap))
		for k, _ := range fieldsMap {
			fieldsKeyMap[k] = nil
		}
		// Automatic data key to table field name mapping.
		var foundKey string
		for dataKey, dataValue := range data {
			if _, ok := fieldsKeyMap[dataKey]; !ok {
				foundKey, _ = gutil.MapPossibleItemByKey(fieldsKeyMap, dataKey)
				if foundKey != "" {
					data[foundKey] = dataValue
					delete(data, dataKey)
				} else if !filter {
					if schema != "" {
						return nil, gerror.Newf(`no column of name "%s" found for table "%s" in schema "%s"`, dataKey, table, schema)
					}
					return nil, gerror.Newf(`no column of name "%s" found for table "%s"`, dataKey, table)
				}
			}
		}
		// Data filtering.
		if filter {
			for dataKey, _ := range data {
				if _, ok := fieldsMap[dataKey]; !ok {
					delete(data, dataKey)
				}
			}
		}
	}
	return data, nil
}

//// filterFields removes all key-value pairs which are not the field of given table.
//func (c *Core) filterFields(schema, table string, data map[string]interface{}) map[string]interface{} {
//	// It must use data copy here to avoid its changing the origin data map.
//	newDataMap := make(map[string]interface{}, len(data))
//	if fields, err := c.DB.TableFields(table, schema); err == nil {
//		for k, v := range data {
//			if _, ok := fields[k]; ok {
//				newDataMap[k] = v
//			}
//		}
//	}
//	return newDataMap
//}
