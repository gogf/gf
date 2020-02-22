// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/os/gtime"

	"github.com/gogf/gf/encoding/gbinary"

	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// convertValue automatically checks and converts field value from database type
// to golang variable type.
func (bs *dbBase) convertValue(fieldValue []byte, fieldType string) interface{} {
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
func (bs *dbBase) filterFields(schema, table string, data map[string]interface{}) map[string]interface{} {
	// It must use data copy here to avoid its changing the origin data map.
	newDataMap := make(map[string]interface{}, len(data))
	if fields, err := bs.db.TableFields(table, schema); err == nil {
		for k, v := range data {
			if _, ok := fields[k]; ok {
				newDataMap[k] = v
			}
		}
	}
	return newDataMap
}

// Tables returns the table name array of current schema.
func (bs *dbBase) Tables(schema ...string) (tables []string, err error) {
	var result Result
	link, err := bs.db.getSlave(schema...)
	if err != nil {
		return nil, err
	}
	result, err = bs.db.doGetAll(link, `SHOW TABLES`)
	if err != nil {
		return
	}
	for _, m := range result {
		for _, v := range m {
			tables = append(tables, v.String())
		}
	}
	return
}

// TableFields retrieves and returns the fields of given table.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the TableField struct has a "Index" field marks its sequence in the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the process restarts.
func (bs *dbBase) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	table = gstr.Trim(table)
	if gstr.Contains(table, " ") {
		panic("function TableFields supports only single table operations")
	}
	checkSchema := bs.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v := bs.cache.GetOrSetFunc(
		fmt.Sprintf(`mysql_table_fields_%s_%s`, table, checkSchema),
		func() interface{} {
			var result Result
			var link *sql.DB
			link, err = bs.db.getSlave(checkSchema)
			if err != nil {
				return nil
			}
			result, err = bs.doGetAll(
				link,
				fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, bs.db.quoteWord(table)),
			)
			if err != nil {
				return nil
			}
			fields = make(map[string]*TableField)
			for i, m := range result {
				fields[m["Field"].String()] = &TableField{
					Index:   i,
					Name:    m["Field"].String(),
					Type:    m["Type"].String(),
					Null:    m["Null"].Bool(),
					Key:     m["Key"].String(),
					Default: m["Default"].Val(),
					Extra:   m["Extra"].String(),
					Comment: m["Comment"].String(),
				}
			}
			return fields
		}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}
