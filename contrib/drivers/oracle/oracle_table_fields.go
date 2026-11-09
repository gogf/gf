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
// A table qualified with its owner, like `SYSTEM.user`, is looked up in the schema of that owner,
// and a quoted table name, like `"user"`, as it is written.
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result       gdb.Result
		link         gdb.Link
		usedSchema   = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		structureSql string
	)
	owner, name, ok := splitTableName(table)
	if !ok {
		name = strings.ToUpper(table)
	}
	owner, name = escapeSqlString(owner), escapeSqlString(name)
	if owner != "" {
		structureSql = fmt.Sprintf(ownerTableFieldsSqlTmp, owner, name, owner, name)
	} else {
		structureSql = fmt.Sprintf(tableFieldsSqlTmp, name, name)
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
// keeping its owner if it is qualified with one and its quotes if it is quoted, so that TableFields
// looks it up in that owner's schema and as it is written. It returns an empty string if the first
// table is qualified or quoted but not a table name that splitTableName accepts.
func (d *Driver) GetTableNameForFields(table string) string {
	if tables := gstr.SplitAndTrim(table, ","); len(tables) > 0 {
		first := gstr.SplitAndTrim(tables[0], " ")[0]
		if _, _, ok := splitTableName(first); ok {
			return first
		}
		if strings.ContainsAny(first, `."`) {
			return ""
		}
	}
	return d.Core.GetTableNameForFields(table)
}

// escapeSqlString escapes `s` for a string literal in sql by doubling its single quotes.
func escapeSqlString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// splitTableName splits `table`, a table name optionally qualified with its owner, like `user`,
// `"user"`, `SYSTEM.user` or `"SYSTEM"."user"`, into the owner, which is empty if `table` is not
// qualified, and the table name as Oracle stores them: an unquoted identifier in upper case, and
// a quoted identifier as it is, without its quote characters. It returns false if `table` is not
// such a name.
func splitTableName(table string) (owner, name string, ok bool) {
	var end int
	if name, end, ok = scanIdentifier(table, 0); !ok {
		return "", "", false
	}
	if end == len(table) {
		return "", name, true
	}
	if table[end] != '.' {
		return "", "", false
	}
	owner = name
	if name, end, ok = scanIdentifier(table, end+1); !ok || end != len(table) {
		return "", "", false
	}
	return owner, name, true
}

// scanIdentifier returns the identifier that begins at index `i` of `s` as Oracle stores it, being
// an unquoted identifier in upper case and a quoted identifier as it is, without its quote
// characters, and the index following it. It returns false if no identifier begins there.
func scanIdentifier(s string, i int) (identifier string, end int, ok bool) {
	if i < len(s) && s[i] == '"' {
		n := strings.IndexByte(s[i+1:], '"')
		if n <= 0 {
			return "", 0, false
		}
		return s[i+1 : i+1+n], i + n + 2, true
	}
	end = i
	for end < len(s) && isIdentifierChar(s[end]) {
		end++
	}
	if end == i {
		return "", 0, false
	}
	return strings.ToUpper(s[i:end]), end, true
}

// isIdentifierChar reports whether `c` can be part of an unquoted identifier: a letter, a digit,
// `_`, `$`, `#`, or a byte of a multibyte character.
func isIdentifierChar(c byte) bool {
	return c == '_' || c == '$' || c == '#' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
