// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"

	_ "github.com/gf-third/mysql"
)

type dbMysql struct {
	*dbBase
}

// Open creates and returns a underlying sql.DB object for mysql.
func (db *dbMysql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true&parseTime=true&loc=Local",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("gf-mysql", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// getChars returns the security char for this type of database.
func (db *dbMysql) getChars() (charLeft string, charRight string) {
	return "`", "`"
}

// handleSqlBeforeExec handles the sql before posts it to database.
func (db *dbMysql) handleSqlBeforeExec(sql string) string {
	return sql
}

// Tables retrieves and returns the tables of current schema.
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

// TableFields retrieves and returns the fields information of specified table of current schema.
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
