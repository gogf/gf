// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/lib/pq"
// 2. It does not support Save/Replace features.
// 3. It does not support LastInsertId.

// Package pgsql implements gdb.Driver, which supports operations for PostgreSql.
package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// Driver is the driver for postgresql database.
type Driver struct {
	*gdb.Core
}

var (
	// tableFieldsMap caches the table information retrieved from database.
	tableFieldsMap = gmap.New(true)
)

func init() {
	if err := gdb.Register(`pgsql`, New()); err != nil {
		panic(err)
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for PostgreSql.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for postgresql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

// Open creates and returns an underlying sql.DB object for pgsql.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "postgres"
	)
	if config.Link != "" {
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`dbname=([\w\.\-]+)+`, "dbname="+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
		if config.Timezone != "" {
			source = fmt.Sprintf("%s timezone=%s", source, config.Timezone)
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
	linkInfo := d.GetConfig().Link
	if linkInfo == "" {
		return ""
	}
	s, _ := gregex.ReplaceString(
		`(.+?)\s*password=(.+)\s*host=(.+)`,
		`$1 password=xxx host=$3`,
		linkInfo,
	)
	return s
}

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return `"`, `"`
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	defer func() {
		newSql, newArgs, err = d.Core.DoFilter(ctx, link, newSql, newArgs)
	}()
	var index int
	// Convert placeholder char '?' to string "$x".
	sql, _ = gregex.ReplaceStringFunc(`\?`, sql, func(s string) string {
		index++
		return fmt.Sprintf(`$%d`, index)
	})
	// Handle pgsql jsonb feature support, which contains place holder char '?'.
	// Refer:
	// https://github.com/gogf/gf/issues/1537
	// https://www.postgresql.org/docs/12/functions-json.html
	sql, _ = gregex.ReplaceStringFuncMatch(`(::jsonb([^\w\d]*)\$\d)`, sql, func(match []string) string {
		return fmt.Sprintf(`::jsonb%s?`, match[2])
	})
	newSql, _ = gregex.ReplaceString(` LIMIT (\d+),\s*(\d+)`, ` LIMIT $2 OFFSET $1`, sql)
	return newSql, args, nil
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	query := "SELECT TABLENAME FROM PG_TABLES WHERE SCHEMANAME = 'public' ORDER BY TABLENAME"
	if len(schema) > 0 && schema[0] != "" {
		query = fmt.Sprintf(
			"SELECT TABLENAME FROM PG_TABLES WHERE SCHEMANAME = '%s' ORDER BY TABLENAME",
			schema[0],
		)
	}
	result, err = d.DoSelect(ctx, link, query)
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
		return nil, gerror.NewCode(
			gcode.CodeInvalidParameter,
			"function TableFields supports only single table operations",
		)
	}
	table, _ = gregex.ReplaceString("\"", "", table)
	useSchema := d.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`pgsql_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result       gdb.Result
				link         gdb.Link
				structureSql = fmt.Sprintf(`
SELECT a.attname AS field, t.typname AS type,a.attnotnull as null,
    (case when d.contype is not null then 'pri' else '' end)  as key
      ,ic.column_default as default_value,b.description as comment
      ,coalesce(character_maximum_length, numeric_precision, -1) as length
      ,numeric_scale as scale
FROM pg_attribute a
         left join pg_class c on a.attrelid = c.oid
         left join pg_constraint d on d.conrelid = c.oid and a.attnum = d.conkey[1]
         left join pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid
         left join  pg_type t ON  a.atttypid = t.oid
         left join information_schema.columns ic on ic.column_name = a.attname and ic.table_name = c.relname
WHERE c.relname = '%s' and a.attisdropped is false and a.attnum > 0
ORDER BY a.attnum`,
					table,
				)
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			structureSql, _ = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(structureSql))
			result, err = d.DoSelect(ctx, link, structureSql)
			if err != nil {
				return nil
			}
			fields = make(map[string]*gdb.TableField)
			for i, m := range result {
				fields[m["field"].String()] = &gdb.TableField{
					Index:   i,
					Name:    m["field"].String(),
					Type:    m["type"].String(),
					Null:    !m["null"].Bool(),
					Key:     m["key"].String(),
					Default: m["default_value"].Val(),
					Comment: m["comment"].String(),
				}
			}
			return fields
		},
	)
	if v != nil {
		fields = v.(map[string]*gdb.TableField)
	}
	return
}

// DoInsert is not supported in pgsql.
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Save operation is not supported by pgsql driver`,
		)

	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by pgsql driver`,
		)

	default:
		return d.Core.DoInsert(ctx, link, table, list, option)
	}
}
