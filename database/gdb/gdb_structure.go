// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"strings"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/os/gtime"

	"github.com/gogf/gf/encoding/gbinary"

	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// convertValue automatically checks and converts field value from database type
// to golang variable type.
func (c *Core) convertValue(fieldValue []byte, fieldType string) interface{} {
	t, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	t = strings.ToLower(t)
	switch t {
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		return fieldValue

	case "int", "tinyint", "small_int", "smallint", "medium_int", "mediumint":
		if gstr.ContainsI(fieldType, "unsigned") {
			gconv.Uint(string(fieldValue))
		}
		return gconv.Int(string(fieldValue))

	case "big_int", "bigint":
		if gstr.ContainsI(fieldType, "unsigned") {
			gconv.Uint64(string(fieldValue))
		}
		return gconv.Int64(string(fieldValue))

	case "float", "double", "decimal":
		return gconv.Float64(string(fieldValue))

	case "bit":
		s := string(fieldValue)
		// mssql is true|false string.
		if strings.EqualFold(s, "true") {
			return 1
		}
		if strings.EqualFold(s, "false") {
			return 0
		}
		return gbinary.BeDecodeToInt64(fieldValue)

	case "bool":
		return gconv.Bool(fieldValue)

	case "date":
		t, _ := gtime.StrToTime(string(fieldValue))
		return t.Format("Y-m-d")

	case "datetime", "timestamp":
		t, _ := gtime.StrToTime(string(fieldValue))
		return t.String()

	default:
		// Auto detect field type, using key match.
		switch {
		case strings.Contains(t, "text") || strings.Contains(t, "char"):
			return string(fieldValue)

		case strings.Contains(t, "float") || strings.Contains(t, "double"):
			return gconv.Float64(string(fieldValue))

		case strings.Contains(t, "bool"):
			return gconv.Bool(string(fieldValue))

		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			return fieldValue

		case strings.Contains(t, "int"):
			return gconv.Int(string(fieldValue))

		case strings.Contains(t, "time"):
			s := string(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t.String()

		case strings.Contains(t, "date"):
			s := string(fieldValue)
			t, err := gtime.StrToTime(s)
			if err != nil {
				return s
			}
			return t.Format("Y-m-d")

		default:
			return string(fieldValue)
		}
	}
}

// filterFields removes all key-value pairs which are not the field of given table.
func (c *Core) filterFields(schema, table string, data map[string]interface{}) map[string]interface{} {
	// It must use data copy here to avoid its changing the origin data map.
	newDataMap := make(map[string]interface{}, len(data))
	if fields, err := c.DB.TableFields(table, schema); err == nil {
		for k, v := range data {
			if _, ok := fields[k]; ok {
				newDataMap[k] = v
			}
		}
	}
	return newDataMap
}
