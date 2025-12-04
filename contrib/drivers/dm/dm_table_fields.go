// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

// escapeSingleQuote escapes single quotes in the string to prevent SQL injection.
// In SQL, single quotes are escaped by doubling them (two single quotes).
func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

const (
	tableFieldsSqlTmp         = `SELECT c.COLUMN_NAME, c.DATA_TYPE, c.DATA_LENGTH, c.DATA_DEFAULT, c.NULLABLE, cc.COMMENTS FROM ALL_TAB_COLUMNS c LEFT JOIN ALL_COL_COMMENTS cc ON c.COLUMN_NAME = cc.COLUMN_NAME AND c.TABLE_NAME = cc.TABLE_NAME AND c.OWNER = cc.OWNER WHERE c.TABLE_NAME = '%s' AND c.OWNER = '%s'`
	tableFieldsPkSqlSchemaTmp = `SELECT COLS.COLUMN_NAME AS PRIMARY_KEY_COLUMN FROM USER_CONSTRAINTS CONS JOIN USER_CONS_COLUMNS COLS ON CONS.CONSTRAINT_NAME = COLS.CONSTRAINT_NAME WHERE CONS.TABLE_NAME = '%s' AND CONS.CONSTRAINT_TYPE = 'P'`
	tableFieldsPkSqlDBATmp    = `SELECT COLS.COLUMN_NAME AS PRIMARY_KEY_COLUMN FROM DBA_CONSTRAINTS CONS JOIN DBA_CONS_COLUMNS COLS ON CONS.CONSTRAINT_NAME = COLS.CONSTRAINT_NAME WHERE CONS.TABLE_NAME = '%s' AND CONS.OWNER = '%s' AND CONS.CONSTRAINT_TYPE = 'P'`
)

// TableFields retrieves and returns the fields' information of specified table of current schema.
func (d *Driver) TableFields(
	ctx context.Context, table string, schema ...string,
) (fields map[string]*gdb.TableField, err error) {
	var (
		result   gdb.Result
		pkResult gdb.Result
		link     gdb.Link
		// When no schema is specified, the configuration item is returned by default
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	// When usedSchema is empty, return the default link
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	// The link has been distinguished and no longer needs to judge the owner
	result, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(
			tableFieldsSqlTmp,
			escapeSingleQuote(strings.ToUpper(table)),
			escapeSingleQuote(strings.ToUpper(d.GetSchema())),
		),
	)
	if err != nil {
		return nil, err
	}
	// Query the primary key field
	pkResult, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(tableFieldsPkSqlSchemaTmp, escapeSingleQuote(strings.ToUpper(table))),
	)
	if err != nil {
		return nil, err
	}
	if pkResult.IsEmpty() {
		pkResult, err = d.DoSelect(
			ctx, link,
			fmt.Sprintf(tableFieldsPkSqlDBATmp, escapeSingleQuote(strings.ToUpper(table)), escapeSingleQuote(strings.ToUpper(d.GetSchema()))),
		)
		if err != nil {
			return nil, err
		}
	}
	fields = make(map[string]*gdb.TableField)
	pkFields := gmap.NewStrStrMap()
	for _, pk := range pkResult {
		pkFields.Set(pk["PRIMARY_KEY_COLUMN"].String(), "PRI")
	}
	for i, m := range result {
		// m[NULLABLE] returns "N" "Y"
		// "N" means not null
		// "Y" means could be null
		var nullable bool
		if m["NULLABLE"].String() != "N" {
			nullable = true
		}

		// Build field type with length/precision
		// For NUMBER(p,s): use DATA_PRECISION and DATA_SCALE
		// For VARCHAR2/CHAR: use DATA_LENGTH
		var (
			fieldType  string
			dataType   = m["DATA_TYPE"].String()
			dataLength = m["DATA_LENGTH"].Int()
		)
		if dataLength > 0 {
			fieldType = fmt.Sprintf("%s(%d)", dataType, dataLength)
		} else {
			fieldType = dataType
		}
		fields[m["COLUMN_NAME"].String()] = &gdb.TableField{
			Index:   i,
			Name:    m["COLUMN_NAME"].String(),
			Type:    fieldType,
			Null:    nullable,
			Default: m["DATA_DEFAULT"].Val(),
			Key:     pkFields.Get(m["COLUMN_NAME"].String()),
			// Extra:   m["Extra"].String(),
			Comment: m["COMMENTS"].String(),
		}
	}
	return fields, nil
}
