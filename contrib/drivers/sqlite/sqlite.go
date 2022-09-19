// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/glebarez/go-sqlite"
// 2. It does not support Save/Replace features.

// Package sqlite implements gdb.Driver, which supports operations for SQLite.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/gogf/gf/v2/encoding/gurl"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Driver is the driver for sqlite database.
type Driver struct {
	*gdb.Core
}

var (
	ErrorSave = gerror.NewCode(gcode.CodeNotSupported, `Save operation is not supported by sqlite driver`)
)

func init() {
	if err := gdb.Register(`sqlite`, New()); err != nil {
		panic(err)
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for SQLite.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for sqlite.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for sqlite.
// https://github.com/glebarez/go-sqlite
func (d *Driver) Open(config gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "sqlite"
	)
	if config.Link != "" {
		// ============================================================================
		// Deprecated from v2.2.0.
		// ============================================================================
		source = config.Link
	} else {
		source = config.Name
	}
	// It searches the source file to locate its absolute path..
	if absolutePath, _ := gfile.Search(source); absolutePath != "" {
		source = absolutePath
	}

	// Multiple PRAGMAs can be specified, e.g.:
	// path/to/some.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)
	if config.Extra != "" {
		var (
			options  string
			extraMap map[string]interface{}
		)
		if extraMap, err = gstr.Parse(config.Extra); err != nil {
			return nil, err
		}
		for k, v := range extraMap {
			if options != "" {
				options += "&"
			}
			options += fmt.Sprintf(`_pragma=%s(%s)`, k, gurl.Encode(gconv.String(v)))
		}
		if len(options) > 1 {
			source += "?" + options
		}
	}

	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}

// FilteredLink retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *Driver) FilteredLink() string {
	return d.GetConfig().Link
}

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	return d.Core.DoFilter(ctx, link, sql, args)
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}

	result, err = d.DoSelect(ctx, link, `SELECT NAME FROM SQLITE_MASTER WHERE TYPE='table' ORDER BY NAME`)
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

// TableFields retrieves and returns the fields' information of specified table of current schema.
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "function TableFields supports only single table operations")
	}
	useSchema := d.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	var (
		result gdb.Result
		link   gdb.Link
	)
	if link, err = d.SlaveLink(useSchema); err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, fmt.Sprintf(`PRAGMA TABLE_INFO(%s)`, table))
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		mKey := ""
		if m["pk"].Bool() {
			mKey = "pri"
		}
		fields[strings.ToLower(m["name"].String())] = &gdb.TableField{
			Index:   i,
			Name:    strings.ToLower(m["name"].String()),
			Type:    strings.ToLower(m["type"].String()),
			Key:     mKey,
			Default: m["dflt_value"].Val(),
			Null:    !m["notnull"].Bool(),
		}
	}
	return fields, nil
}

// DoInsert is not supported in sqlite.
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return nil, ErrorSave
	case gdb.InsertOptionIgnore, gdb.InsertOptionReplace:
		var (
			keys           []string      // Field names.
			values         []string      // Value holder string array, like: (?,?,?)
			params         []interface{} // Values that will be committed to underlying database driver.
			onDuplicateStr string        // onDuplicateStr is used in "ON DUPLICATE KEY UPDATE" statement.
		)
		// Handle the field names and placeholders.
		for k := range list[0] {
			keys = append(keys, k)
		}
		// Prepare the batch result pointer.
		var (
			charL, charR = d.GetChars()
			batchResult  = new(gdb.SqlResult)
			keysStr      = charL + strings.Join(keys, charR+","+charL) + charR
			operation    = "INSERT OR IGNORE"
		)

		if option.InsertOption == gdb.InsertOptionReplace {
			operation = "INSERT OR REPLACE"
		}
		var (
			listLength  = len(list)
			valueHolder = make([]string, 0)
		)
		for i := 0; i < listLength; i++ {
			values = values[:0]
			// Note that the map type is unordered,
			// so it should use slice+key to retrieve the value.
			for _, k := range keys {
				if s, ok := list[i][k].(gdb.Raw); ok {
					values = append(values, gconv.String(s))
				} else {
					values = append(values, "?")
					params = append(params, list[i][k])
				}
			}
			valueHolder = append(valueHolder, "("+gstr.Join(values, ",")+")")
			// Batch package checks: It meets the batch number, or it is the last element.
			if len(valueHolder) == option.BatchCount || (i == listLength-1 && len(valueHolder) > 0) {
				var (
					stdSqlResult sql.Result
					affectedRows int64
				)
				stdSqlResult, err = d.DoExec(ctx, link, fmt.Sprintf(
					"%s INTO %s(%s) VALUES%s %s",
					operation, d.QuotePrefixTableName(table), keysStr,
					gstr.Join(valueHolder, ","),
					onDuplicateStr,
				), params...)
				if err != nil {
					return stdSqlResult, err
				}
				if affectedRows, err = stdSqlResult.RowsAffected(); err != nil {
					err = gerror.WrapCode(gcode.CodeDbOperationError, err, `sql.Result.RowsAffected failed`)
					return stdSqlResult, err
				} else {
					batchResult.Result = stdSqlResult
					batchResult.Affected += affectedRows
				}
				params = params[:0]
				valueHolder = valueHolder[:0]
			}
		}
		return batchResult, nil
	default:
		return d.Core.DoInsert(ctx, link, table, list, option)
	}
}
