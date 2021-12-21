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

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// DriverPgsql is the driver for postgresql database.
type DriverPgsql struct {
	*Core
}

// New creates and returns a database object for postgresql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverPgsql) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverPgsql{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for pgsql.
func (d *DriverPgsql) Open(config *ConfigNode) (db *sql.DB, err error) {
	var (
		source string
		driver = "postgres"
	)
	if config.Link != "" {
		source = config.Link
	} else {
		source = fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
		if config.Timezone != "" {
			source = fmt.Sprintf("%s timezone=%s", source, config.Timezone)
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
func (d *DriverPgsql) FilteredLink() string {
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
func (d *DriverPgsql) GetChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// DoCommit deals with the sql string before commits it to underlying sql driver.
func (d *DriverPgsql) DoCommit(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	defer func() {
		newSql, newArgs, err = d.Core.DoCommit(ctx, link, newSql, newArgs)
	}()

	var index int
	// Convert placeholder char '?' to string "$x".
	sql, _ = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf("$%d", index)
	})
	newSql, _ = gregex.ReplaceString(` LIMIT (\d+),\s*(\d+)`, ` LIMIT $2 OFFSET $1`, sql)
	return newSql, args, nil
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverPgsql) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	query := "SELECT TABLENAME FROM PG_TABLES WHERE SCHEMANAME = 'public' ORDER BY TABLENAME"
	if len(schema) > 0 && schema[0] != "" {
		query = fmt.Sprintf("SELECT TABLENAME FROM PG_TABLES WHERE SCHEMANAME = '%s' ORDER BY TABLENAME", schema[0])
	}
	result, err = d.DoGetAll(ctx, link, query)
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
func (d *DriverPgsql) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "function TableFields supports only single table operations")
	}
	table, _ = gregex.ReplaceString("\"", "", table)
	useSchema := d.db.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`pgsql_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result       Result
				link         Link
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
WHERE c.relname = '%s' and a.attnum > 0
ORDER BY a.attnum`,
					strings.ToLower(table),
				)
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			structureSql, _ = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(structureSql))
			result, err = d.DoGetAll(ctx, link, structureSql)
			if err != nil {
				return nil
			}
			fields = make(map[string]*TableField)
			for i, m := range result {
				fields[m["field"].String()] = &TableField{
					Index:   i,
					Name:    m["field"].String(),
					Type:    m["type"].String(),
					Null:    m["null"].Bool(),
					Key:     m["key"].String(),
					Default: m["default_value"].Val(),
					Comment: m["comment"].String(),
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

// DoInsert is not supported in pgsql.
func (d *DriverPgsql) DoInsert(ctx context.Context, link Link, table string, list List, option DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case insertOptionSave:
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Save operation is not supported by pgsql driver`)

	case insertOptionReplace:
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Replace operation is not supported by pgsql driver`)

	default:
		return d.Core.DoInsert(ctx, link, table, list, option)
	}
}
