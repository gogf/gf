// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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
	"context"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/errors/gcode"
	"strings"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
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
	if config.Link != "" {
		source = config.Link
	} else {
		source = config.Name
	}
	// It searches the source file to locate its absolute path..
	if absolutePath, _ := gfile.Search(source); absolutePath != "" {
		source = absolutePath
	}
	intlog.Printf(d.GetCtx(), "Open: %s", source)
	if db, err := sql.Open("sqlite3", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// FilteredLink retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *DriverSqlite) FilteredLink() string {
	return d.GetConfig().Link
}

// GetChars returns the security char for this type of database.
func (d *DriverSqlite) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// DoCommit deals with the sql string before commits it to underlying sql driver.
func (d *DriverSqlite) DoCommit(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	return d.Core.DoCommit(ctx, link, sql, args)
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverSqlite) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}

	result, err = d.DoGetAll(ctx, link, `SELECT NAME FROM SQLITE_MASTER WHERE TYPE='table' ORDER BY NAME`)
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
// Also see DriverMysql.TableFields.
func (d *DriverSqlite) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "function TableFields supports only single table operations")
	}
	useSchema := d.db.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	tableFieldsCacheKey := fmt.Sprintf(
		`sqlite_table_fields_%s_%s@group:%s`,
		table, useSchema, d.GetGroup(),
	)
	v := tableFieldsMap.GetOrSetFuncLock(tableFieldsCacheKey, func() interface{} {
		var (
			result    Result
			link, err = d.SlaveLink(useSchema)
		)
		if err != nil {
			return nil
		}
		result, err = d.DoGetAll(ctx, link, fmt.Sprintf(`PRAGMA TABLE_INFO(%s)`, table))
		if err != nil {
			return nil
		}
		fields = make(map[string]*TableField)
		for i, m := range result {
			fields[strings.ToLower(m["name"].String())] = &TableField{
				Index: i,
				Name:  strings.ToLower(m["name"].String()),
				Type:  strings.ToLower(m["type"].String()),
			}
		}
		return fields
	})
	if v != nil {
		fields = v.(map[string]*TableField)
	}
	return
}

// DoInsert is not supported in sqlite.
func (d *DriverSqlite) DoInsert(ctx context.Context, link Link, table string, list List, option DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case insertOptionSave:
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Save operation is not supported by sqlite driver`)

	case insertOptionReplace:
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Replace operation is not supported by sqlite driver`)

	default:
		return d.Core.DoInsert(ctx, link, table, list, option)
	}
}

//ExpandFields 获取扩展列信息
func (d *DriverSqlite) ExpandFields(ctx context.Context, bizCode, bizType string, params ...string) (columns []*ExpandField, err error) {
	return nil, nil
}
