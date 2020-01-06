// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/lib/pq"
// 2. It does not support Save/Replace features.
// 3. It does not support LastInsertId.

package gdb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/text/gregex"
)

type dbPgsql struct {
	*dbBase
}

func (db *dbPgsql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
	}
	if db, err := sql.Open("postgres", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

func (db *dbPgsql) setSchema(sqlDb *sql.DB, schema string) error {
	_, err := sqlDb.Exec("SET search_path TO " + schema)
	return err
}

func (db *dbPgsql) getChars() (charLeft string, charRight string) {
	return "\"", "\""
}

func (db *dbPgsql) handleSqlBeforeExec(query string) string {
	index := 0
	query, _ = gregex.ReplaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf("$%d", index)
	})
	query, _ = gregex.ReplaceString(` LIMIT (\d+),\s*(\d+)`, ` LIMIT $1 OFFSET $2`, query)
	return query
}

// 返回当前数据库所有的数据表名称
// TODO
func (db *dbPgsql) Tables() (tables []string, err error) {
	return
}

// 获得指定表表的数据结构，构造成map哈希表返回，其中键名为表字段名称，键值为字段数据结构.
func (db *dbPgsql) TableFields(table string) (fields map[string]*TableField, err error) {
	table, _ = gregex.ReplaceString("\"", "", table)
	v := db.cache.GetOrSetFunc("pgsql_table_fields_"+table, func() interface{} {
		result := (Result)(nil)
		result, err = db.GetAll(fmt.Sprintf(`
		SELECT a.attname AS field, t.typname AS type FROM pg_class c, pg_attribute a 
        LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid,pg_type t
        WHERE c.relname = '%s' and a.attnum > 0 and a.attrelid = c.oid and a.atttypid = t.oid ORDER BY a.attnum`, strings.ToLower(table)))
		if err != nil {
			return nil
		}

		fields = make(map[string]*TableField)
		for i, m := range result {
			fields[m["field"].String()] = &TableField{
				Index: i,
				Name:  m["field"].String(),
				Type:  m["type"].String(),
			}
		}
		return fields
	}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}
