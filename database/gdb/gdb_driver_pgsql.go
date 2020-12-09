// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
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
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"
	"strings"

	"github.com/gogf/gf/text/gregex"
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
func (d *DriverPgsql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("postgres", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// GetChars returns the security char for this type of database.
func (d *DriverPgsql) GetChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// HandleSqlBeforeCommit deals with the sql string before commits it to underlying sql driver.
func (d *DriverPgsql) HandleSqlBeforeCommit(link Link, sql string, args []interface{}) (string, []interface{}) {
	var index int
	// Convert place holder char '?' to string "$x".
	sql, _ = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf("$%d", index)
	})
	sql, _ = gregex.ReplaceString(` LIMIT (\d+),\s*(\d+)`, ` LIMIT $2 OFFSET $1`, sql)
	return sql, args
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverPgsql) Tables(schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.DB.GetSlave(schema...)
	if err != nil {
		return nil, err
	}
	query := "SELECT TABLENAME FROM PG_TABLES WHERE SCHEMANAME = 'public' ORDER BY TABLENAME"
	if len(schema) > 0 && schema[0] != "" {
		query = fmt.Sprintf("SELECT TABLENAME FROM PG_TABLES WHERE SCHEMANAME = '%s' ORDER BY TABLENAME", schema[0])
	}
	result, err = d.DB.DoGetAll(link, query)
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
func (d *DriverPgsql) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.New("function TableFields supports only single table operations")
	}
	table, _ = gregex.ReplaceString("\"", "", table)
	checkSchema := d.DB.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v, _ := internalCache.GetOrSetFunc(
		fmt.Sprintf(`pgsql_table_fields_%s_%s@group:%s`, table, checkSchema, d.GetGroup()),
		func() (interface{}, error) {
			var (
				result Result
				link   *sql.DB
			)
			link, err = d.DB.GetSlave(checkSchema)
			if err != nil {
				return nil, err
			}
			structureSql := fmt.Sprintf(`
SELECT a.attname AS field, t.typname AS type FROM pg_class c, pg_attribute a 
LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid,pg_type t
WHERE c.relname = '%s' and a.attnum > 0 and a.attrelid = c.oid and a.atttypid = t.oid 
ORDER BY a.attnum`,
				strings.ToLower(table),
			)
			structureSql, _ = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(structureSql))
			result, err = d.DB.DoGetAll(link, structureSql)
			if err != nil {
				return nil, err
			}

			fields = make(map[string]*TableField)
			for i, m := range result {
				fields[m["field"].String()] = &TableField{
					Index: i,
					Name:  m["field"].String(),
					Type:  m["type"].String(),
				}
			}
			return fields, nil
		}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}
