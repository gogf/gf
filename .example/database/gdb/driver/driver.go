// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package driver

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"

	_ "github.com/gf-third/mysql"
)

type MyDriver struct {
	*gdb.Core
}

// Open creates and returns a underlying sql.DB object for mysql.
func (d *MyDriver) Open(config *gdb.ConfigNode) (*sql.DB, error) {
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
func (d *MyDriver) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// handleSqlBeforeExec handles the sql before posts it to database.
func (d *MyDriver) HandleSqlBeforeExec(sql string) string {
	return sql
}

// Tables retrieves and returns the tables of current schema.
func (d *MyDriver) Tables(schema ...string) (tables []string, err error) {
	var result gdb.Result
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

// gdb.TableFields retrieves and returns the fields information of specified table of current schema.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the gdb.TableField struct has a "Index" field marks its sequence in the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the process restarts.
func (d *MyDriver) TableFields(table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	table = gstr.Trim(table)
	if gstr.Contains(table, " ") {
		panic("function gdb.TableFields supports only single table operations")
	}
	checkSchema := d.DB.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v := d.DB.GetCache().GetOrSetFunc(
		fmt.Sprintf(`mysql_table_fields_%s_%s`, table, checkSchema),
		func() interface{} {
			var result gdb.Result
			var link *sql.DB
			link, err = d.DB.GetSlave(checkSchema)
			if err != nil {
				return nil
			}
			result, err = d.DB.DoGetAll(
				link,
				fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, d.DB.QuoteWord(table)),
			)
			if err != nil {
				return nil
			}
			fields = make(map[string]*gdb.TableField)
			for i, m := range result {
				fields[m["Field"].String()] = &gdb.TableField{
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
		fields = v.(map[string]*gdb.TableField)
	}
	return
}
