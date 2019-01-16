// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gdb

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gregex"
    "strings"
)

/*
// 同步数据库表结构到内存中
func (bs *dbBase) syncTableStructure() {
    bs.tables = make(map[string]map[string]string)
    for _, table := range bs.db.getTables() {
        bs.tables[table], _ = bs.db.getTableFields(table)
    }
}
*/

// 字段类型转换，将数据库字段类型转换为golang变量类型
func (bs *dbBase) convertValue(fieldValue interface{}, fieldType string) interface{} {
    t, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
    t     = strings.ToLower(t)
    switch t {
    case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
        return gconv.Bytes(fieldValue)

    case "bit", "int", "tinyint", "small_int", "medium_int":
        return gconv.Int(fieldValue)

    case "big_int":
        return gconv.Int64(fieldValue)

    case "float", "double", "decimal":
        return gconv.Float64(fieldValue)

    case "bool":
        return gconv.Bool(fieldValue)

    default:
        // 自动识别类型, 以便默认支持更多数据库类型
        switch {
            case strings.Contains(t, "int"):
                return gconv.Int(fieldValue)

            case strings.Contains(t, "text") || strings.Contains(t, "char"):
                return gconv.String(fieldValue)

            case strings.Contains(t, "float") || strings.Contains(t, "double"):
                return gconv.Float64(fieldValue)

            case strings.Contains(t, "bool"):
                return gconv.Bool(fieldValue)

            case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
                return gconv.Bytes(fieldValue)

            default:
                return gconv.String(fieldValue)
        }
    }
}

// 将map的数据按照fields进行过滤，只保留与表字段同名的数据
func (bs *dbBase) filterFields(table string, data map[string]interface{}) map[string]interface{} {
    if fields, err := bs.db.getTableFields(table); err == nil {
        for k, _ := range data {
            if _, ok := fields[k]; !ok {
                delete(data, k)
            }
        }
    }
    return data
}

// 获得指定表表的数据结构，构造成map哈希表返回，其中键名为表字段名称，键值暂无用途(默认为字段数据类型).
func (bs *dbBase) getTableFields(table string) (fields map[string]string, err error) {
    // 缓存不存在时会查询数据表结构，缓存后不过期，直至程序重启(重新部署)
    v := bs.cache.GetOrSetFunc("table_fields_" + table, func() interface{} {
        result       := (Result)(nil)
        charl, charr := bs.db.getChars()
        result, err   = bs.GetAll(fmt.Sprintf(`SHOW COLUMNS FROM %s%s%s`, charl, table, charr))
        if err != nil {
            return nil
        }
        fields = make(map[string]string)
        for _, m := range result {
            fields[m["Field"].String()] = m["Type"].String()
        }
        return fields
    }, 0)
    if err == nil {
        fields = v.(map[string]string)
    }
    return
}

/*
// 获取当前数据库所有的表结构
func (bs *dbBase) getTables() []string {
    if result, _ := bs.GetAll(`SHOW TABLES`); result != nil {
        array := make([]string, len(result))
        for i, m := range result {
            for _, v := range m {
                array[i] = v.String()
                break
            }
        }
        return array
    }
    return nil
}
*/