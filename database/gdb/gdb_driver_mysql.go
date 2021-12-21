// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
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

// Open creates and returns an underlying sql.DB object for mysql.
// Note that it converts time.Time argument to local timezone in default.
func (d *DriverMysql) Open(config *ConfigNode) (db *sql.DB, err error) {
	var (
		source string
		driver = "mysql"
	)
	if config.Link != "" {
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`/([\w\.\-]+)+`, "/"+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
		if config.Timezone != "" {
			source = fmt.Sprintf("%s&loc=%s", source, url.QueryEscape(config.Timezone))
		}
	}
	intlog.Printf(d.GetCtx(), "Open: %s", source)
	if db, err = sql.Open(driver, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, driver, source,
		)
		return nil, err
	}
	return
}

// FilteredLink retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *DriverMysql) FilteredLink() string {
	linkInfo := d.GetConfig().Link
	if linkInfo == "" {
		return ""
	}
	s, _ := gregex.ReplaceString(
		`(.+?):(.+)@tcp(.+)`,
		`$1:xxx@tcp$3`,
		linkInfo,
	)
	return s
}

// GetChars returns the security char for this type of database.
func (d *DriverMysql) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// DoCommit handles the sql before posts it to database.
func (d *DriverMysql) DoCommit(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	return d.Core.DoCommit(ctx, link, sql, args)
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverMysql) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	result, err = d.DoGetAll(ctx, link, `SHOW TABLES`)
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

// TableFields retrieves and returns the fields' information of specified table of current
// schema.
//
// The parameter `link` is optional, if given nil it automatically retrieves a raw sql connection
// as its link to proceed necessary sql query.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the TableField struct has a "Index" field marks its sequence in
// the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the
// process restarts.
func (d *DriverMysql) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "function TableFields supports only single table operations")
	}
	useSchema := d.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`mysql_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result Result
				link   Link
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			result, err = d.DoGetAll(ctx, link, fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, d.QuoteWord(table)))
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
		},
	)
	if v != nil {
		fields = v.(map[string]*TableField)
	}
	return
}
