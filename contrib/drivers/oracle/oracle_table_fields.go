// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

// tableFieldsColumnsSql selects the name, type and nullability of the columns aliased as `c`.
const tableFieldsColumnsSql = `
    c.COLUMN_NAME AS FIELD, 
    CASE   
    WHEN (c.DATA_TYPE='NUMBER' AND NVL(c.DATA_SCALE,0)=0) THEN 'INT'||'('||c.DATA_PRECISION||','||c.DATA_SCALE||')'
    WHEN (c.DATA_TYPE='NUMBER' AND NVL(c.DATA_SCALE,0)>0) THEN 'FLOAT'||'('||c.DATA_PRECISION||','||c.DATA_SCALE||')'
    WHEN c.DATA_TYPE='FLOAT' THEN c.DATA_TYPE||'('||c.DATA_PRECISION||','||c.DATA_SCALE||')' 
    ELSE c.DATA_TYPE||'('||c.DATA_LENGTH||')' END AS TYPE,
    c.NULLABLE,`

var (
	tableFieldsSqlTmp = `
SELECT ` + tableFieldsColumnsSql + `
    CASE WHEN pk.COLUMN_NAME IS NOT NULL THEN 'PRI' ELSE '' END AS KEY
FROM USER_TAB_COLUMNS c
LEFT JOIN (
    SELECT cols.COLUMN_NAME 
    FROM USER_CONSTRAINTS cons 
    JOIN USER_CONS_COLUMNS cols ON cons.CONSTRAINT_NAME = cols.CONSTRAINT_NAME 
    WHERE cons.TABLE_NAME = '%s' AND cons.CONSTRAINT_TYPE = 'P'
) pk ON c.COLUMN_NAME = pk.COLUMN_NAME
WHERE c.TABLE_NAME = '%s' 
ORDER BY c.COLUMN_ID
`
	ownerTableFieldsSqlTmp = `
SELECT ` + tableFieldsColumnsSql + `
    CASE WHEN c.COLUMN_NAME IN (
        SELECT cols.COLUMN_NAME 
        FROM ALL_CONSTRAINTS cons 
        JOIN ALL_CONS_COLUMNS cols ON cons.OWNER = cols.OWNER AND cons.CONSTRAINT_NAME = cols.CONSTRAINT_NAME 
        WHERE cons.OWNER = '%s' AND cons.TABLE_NAME = '%s' AND cons.CONSTRAINT_TYPE = 'P'
    ) THEN 'PRI' ELSE '' END AS KEY
FROM ALL_TAB_COLUMNS c
WHERE c.OWNER = '%s' AND c.TABLE_NAME = '%s' 
ORDER BY c.COLUMN_ID
`
)

func init() {
	var err error
	tableFieldsSqlTmp, err = gdb.FormatMultiLineSqlToSingle(tableFieldsSqlTmp)
	if err != nil {
		panic(err)
	}
	ownerTableFieldsSqlTmp, err = gdb.FormatMultiLineSqlToSingle(ownerTableFieldsSqlTmp)
	if err != nil {
		panic(err)
	}
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
// A table qualified with its owner, like `SYSTEM.user`, is looked up in the schema of that owner.
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result       gdb.Result
		link         gdb.Link
		usedSchema   = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		structureSql string
	)
	if owner, name, ok := splitOwnerTable(table); ok {
		owner, name = escapeSqlString(owner), escapeSqlString(name)
		structureSql = fmt.Sprintf(ownerTableFieldsSqlTmp, owner, name, owner, name)
	} else {
		upperTable := escapeSqlString(strings.ToUpper(table))
		structureSql = fmt.Sprintf(tableFieldsSqlTmp, upperTable, upperTable)
	}
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, structureSql)
	if err != nil {
		return nil, err
	}

	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		isNull := false
		if m["NULLABLE"].String() == "Y" {
			isNull = true
		}

		fields[m["FIELD"].String()] = &gdb.TableField{
			Index: i,
			Name:  m["FIELD"].String(),
			Type:  m["TYPE"].String(),
			Null:  isNull,
			Key:   m["KEY"].String(),
		}
	}
	return fields, nil
}

// GetTableNameForFields returns the first table in the table expression `table` without its alias,
// keeping its owner if it is qualified with one, so that TableFields looks it up in that owner's schema.
func (d *Driver) GetTableNameForFields(table string) string {
	name := d.Core.GetTableNameForFields(table)
	if name == "" {
		return ""
	}
	if first := gstr.SplitAndTrim(gstr.SplitAndTrim(table, ",")[0], " ")[0]; strings.Count(first, ".") == 1 {
		return first
	}
	return name
}

// escapeSqlString escapes `s` for a string literal in sql by doubling its single quotes.
func escapeSqlString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// splitOwnerTable splits `table` qualified with its owner, like `SYSTEM.user` or `"SYSTEM"."user"`,
// into the owner and the table name, in upper case and without quote characters. It returns false
// if `table` is not qualified.
func splitOwnerTable(table string) (owner, name string, ok bool) {
	if owner, name, ok = strings.Cut(table, "."); !ok {
		return "", "", false
	}
	return strings.ToUpper(strings.Trim(owner, `"`)), strings.ToUpper(strings.Trim(name, `"`)), true
}
