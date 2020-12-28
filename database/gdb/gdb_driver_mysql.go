// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"

	_ "github.com/go-sql-driver/mysql"
)

// DriverMysql is the driver for mysql database.
type DriverMysql struct {
	*Core
}

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverMysql) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverMysql{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for mysql.
// Note that it converts time.Time argument to local timezone in default.
func (d *DriverMysql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`/([\w\.\-]+)+`, "/"+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true&parseTime=true",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("mysql", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// GetChars returns the security char for this type of database.
func (d *DriverMysql) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// HandleSqlBeforeCommit handles the sql before posts it to database.
func (d *DriverMysql) HandleSqlBeforeCommit(link Link, sql string, args []interface{}) (string, []interface{}) {
	return sql, args
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverMysql) Tables(schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.DB.GetSlave(schema...)
	if err != nil {
		return nil, err
	}
	result, err = d.DB.DoGetAll(link, `SHOW TABLES`)
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

// TableFields retrieves and returns the fields information of specified table of current
// schema.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the TableField struct has a "Index" field marks its sequence in
// the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the
// process restarts.
func (d *DriverMysql) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.New("function TableFields supports only single table operations")
	}
	checkSchema := d.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v, _ := internalCache.GetOrSetFunc(
		fmt.Sprintf(`mysql_table_fields_%s_%s@group:%s`, table, checkSchema, d.GetGroup()),
		func() (interface{}, error) {
			var (
				result Result
				link   *sql.DB
			)
			link, err = d.DB.GetSlave(checkSchema)
			if err != nil {
				return nil, err
			}
			result, err = d.DB.DoGetAll(
				link,
				fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, d.DB.QuoteWord(table)),
			)
			if err != nil {
				return nil, err
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
			return fields, nil
		}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}
