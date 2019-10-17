// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/os/gtime"

	"github.com/gogf/gf/encoding/gbinary"

	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// 字段类型转换，将数据库字段类型转换为golang变量类型
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
		// 这里的字符串判断是为兼容不同的数据库类型，如: mssql
		if strings.EqualFold(s, "true") {
			return 1
		}
		if strings.EqualFold(s, "false") {
			return 0
		}
		return gbinary.BeDecodeToInt64(fieldValue)

	case "bool":
		return gconv.Bool(fieldValue)

	case "datetime":
		t, _ := gtime.StrToTime(string(fieldValue))
		return t.String()

	default:
		// 自动识别类型, 以便默认支持更多数据库类型
		switch {
		case strings.Contains(t, "int"):
			return gconv.Int(string(fieldValue))

		case strings.Contains(t, "text") || strings.Contains(t, "char"):
			return string(fieldValue)

		case strings.Contains(t, "float") || strings.Contains(t, "double"):
			return gconv.Float64(string(fieldValue))

		case strings.Contains(t, "bool"):
			return gconv.Bool(string(fieldValue))

		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			return fieldValue

		default:
			return string(fieldValue)
		}
	}
}

// 将map的数据按照fields进行过滤，只保留与表字段同名的数据
func (bs *dbBase) filterFields(table string, data map[string]interface{}) map[string]interface{} {
	// Must use data copy avoiding change the origin data map.
	newDataMap := make(map[string]interface{}, len(data))
	if fields, err := bs.db.TableFields(table); err == nil {
		for k, v := range data {
			if _, ok := fields[k]; ok {
				newDataMap[k] = v
			}
		}
	}
	return newDataMap
}

// 返回当前数据库所有的数据表名称
func (bs *dbBase) Tables() (tables []string, err error) {
	result := (Result)(nil)
	result, err = bs.GetAll(`SHOW TABLES`)
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

// 获得指定表表的数据结构，构造成map哈希表返回，其中键名为表字段名称，键值为字段数据结构.
func (bs *dbBase) TableFields(table string) (fields map[string]*TableField, err error) {
	// 缓存不存在时会查询数据表结构，缓存后不过期，直至程序重启(重新部署)
	v := bs.cache.GetOrSetFunc("table_fields_"+table, func() interface{} {
		result := (Result)(nil)
		result, err = bs.GetAll(fmt.Sprintf(`SHOW COLUMNS FROM %s`, bs.db.quoteWord(table)))
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
			}
		}
		return fields
	}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}
