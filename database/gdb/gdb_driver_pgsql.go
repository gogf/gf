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

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"

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

// FilteredLinkInfo retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *DriverPgsql) FilteredLinkInfo() string {
	linkInfo := d.GetConfig().LinkInfo
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

// HandleSqlBeforeCommit deals with the sql string before commits it to underlying sql driver.
func (d *DriverPgsql) HandleSqlBeforeCommit(ctx context.Context, link Link, sql string, args []interface{}) (string, []interface{}) {
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

// TableFields retrieves and returns the fields information of specified table of current schema.
//
// Also see DriverMysql.TableFields.
func (d *DriverPgsql) TableFields(ctx context.Context, link Link, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.New("function TableFields supports only single table operations")
	}
	table, _ = gregex.ReplaceString("\"", "", table)
	useSchema := d.db.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	tableFieldsCacheKey := fmt.Sprintf(
		`pgsql_table_fields_%s_%s@group:%s`,
		table, useSchema, d.GetGroup(),
	)
	v := tableFieldsMap.GetOrSetFuncLock(tableFieldsCacheKey, func() interface{} {
		var (
			result       Result
			structureSql = fmt.Sprintf(`
SELECT a.attname AS field, t.typname AS type FROM pg_class c, pg_attribute a 
LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid,pg_type t
WHERE c.relname = '%s' and a.attnum > 0 and a.attrelid = c.oid and a.atttypid = t.oid 
ORDER BY a.attnum`,
				strings.ToLower(table),
			)
		)
		structureSql, _ = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(structureSql))
		if link == nil {
			link, err = d.SlaveLink(useSchema)
			if err != nil {
				return nil
			}
		}
		result, err = d.DoGetAll(ctx, link, structureSql)
		if err != nil {
			return nil
		}
		fields = make(map[string]*TableField)
		for i, m := range result {
			fields[m["field"].String()] = &TableField{
				Index: i,
				Name:  m["field"].String(),
				Type:  m["type"].String(),
			}
		}
		return fields
	})
	if v != nil {
		fields = v.(map[string]*TableField)
	}
	return
}
