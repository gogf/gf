// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/mattn/go-sqlite3"
// 2. It does not support Save/Replace features.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"strings"
)

// DriverSqlite is the driver for sqlite database.
type DriverSqlite struct {
	*Core
}

// New creates and returns a database object for sqlite.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverSqlite) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverSqlite{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for sqlite.
func (d *DriverSqlite) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = config.Name
	}
	// It searches the source file to locate its absolute path..
	if absolutePath, _ := gfile.Search(source); absolutePath != "" {
		source = absolutePath
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("sqlite3", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// GetChars returns the security char for this type of database.
func (d *DriverSqlite) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// HandleSqlBeforeCommit deals with the sql string before commits it to underlying sql driver.
// TODO 需要增加对Save方法的支持，可使用正则来实现替换，
// TODO 将ON DUPLICATE KEY UPDATE触发器修改为两条SQL语句(INSERT OR IGNORE & UPDATE)
func (d *DriverSqlite) HandleSqlBeforeCommit(link Link, sql string, args []interface{}) (string, []interface{}) {
	return sql, args
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverSqlite) Tables(schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.DB.GetSlave(schema...)
	if err != nil {
		return nil, err
	}

	result, err = d.DB.DoGetAll(link, `SELECT NAME FROM SQLITE_MASTER WHERE TYPE='table' ORDER BY NAME`)
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
func (d *DriverSqlite) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.New("function TableFields supports only single table operations")
	}
	checkSchema := d.DB.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v, _ := internalCache.GetOrSetFunc(
		fmt.Sprintf(`sqlite_table_fields_%s_%s@group:%s`, table, checkSchema, d.GetGroup()),
		func() (interface{}, error) {
			var (
				result Result
				link   *sql.DB
			)
			link, err = d.DB.GetSlave(checkSchema)
			if err != nil {
				return nil, err
			}
			result, err = d.DB.DoGetAll(link, fmt.Sprintf(`PRAGMA TABLE_INFO(%s)`, table))
			if err != nil {
				return nil, err
			}
			fields = make(map[string]*TableField)
			for i, m := range result {
				fields[strings.ToLower(m["name"].String())] = &TableField{
					Index: i,
					Name:  strings.ToLower(m["name"].String()),
					Type:  strings.ToLower(m["type"].String()),
				}
			}
			return fields, nil
		}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}
