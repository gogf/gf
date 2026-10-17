// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	tableFieldsSqlTmp = `
SELECT 
	c.name AS Field,
	CASE 
		WHEN t.name IN ('datetime', 'datetime2', 'smalldatetime', 'date', 'time', 'text', 'ntext', 'image', 'xml') THEN t.name
		WHEN t.name IN ('decimal', 'numeric') THEN t.name + '(' + CAST(c.precision AS varchar(20)) + ',' + CAST(c.scale AS varchar(20)) + ')'
		WHEN t.name IN ('char', 'varchar', 'binary', 'varbinary') THEN t.name + '(' + CASE WHEN c.max_length = -1 THEN 'max' ELSE CAST(c.max_length AS varchar(20)) END + ')'
		WHEN t.name IN ('nchar', 'nvarchar') THEN t.name + '(' + CASE WHEN c.max_length = -1 THEN 'max' ELSE CAST(c.max_length/2 AS varchar(20)) END + ')'
		ELSE t.name
	END AS Type,
	CASE WHEN c.is_nullable = 1 THEN 'YES' ELSE 'NO' END AS [Null],
	CASE WHEN pk.column_id IS NOT NULL THEN 'PRI' ELSE '' END AS [Key],
	CASE WHEN c.is_identity = 1 THEN 'IDENTITY' ELSE '' END AS Extra,
	ISNULL(dc.definition, '') AS [Default],
	ISNULL(CAST(ep.value AS nvarchar(max)), '') AS [Comment]
FROM sys.columns c
INNER JOIN sys.objects o ON c.object_id = o.object_id AND o.type IN ('U', 'V') AND o.is_ms_shipped = 0
INNER JOIN sys.types t ON c.user_type_id = t.user_type_id
LEFT JOIN sys.default_constraints dc ON c.default_object_id = dc.object_id
LEFT JOIN sys.extended_properties ep ON c.object_id = ep.major_id AND c.column_id = ep.minor_id AND ep.name = 'MS_Description'
LEFT JOIN (
	SELECT ic.object_id, ic.column_id
	FROM sys.index_columns ic
	INNER JOIN sys.indexes i ON ic.object_id = i.object_id AND ic.index_id = i.index_id
	WHERE i.is_primary_key = 1
) pk ON c.object_id = pk.object_id AND c.column_id = pk.column_id
WHERE c.object_id = COALESCE(
	(
		SELECT OBJECT_ID(s.base_object_name)
		FROM sys.synonyms s
		WHERE s.object_id = OBJECT_ID(COALESCE(QUOTENAME(NULLIF(?, N'')) + N'.', N'') + QUOTENAME(?))
		AND PARSENAME(s.base_object_name, 4) IS NULL
		AND COALESCE(PARSENAME(s.base_object_name, 3), DB_NAME()) = DB_NAME()
	),
	OBJECT_ID(COALESCE(QUOTENAME(NULLIF(?, N'')) + N'.', N'') + QUOTENAME(?))
)
ORDER BY c.column_id
`
)

func init() {
	var err error
	tableFieldsSqlTmp, err = gdb.FormatMultiLineSqlToSingle(tableFieldsSqlTmp)
	if err != nil {
		panic(err)
	}
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
// The table is resolved as SQL Server resolves it in a statement: a table name without schema in
// the default schema of the user, a table qualified with its schema, like `dbo.user`, in that
// schema, and a table qualified with its database too, like `test.dbo.user`, in that database. A
// view is resolved like a table, and a synonym as the table or view of the same database it stands for.
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result                  gdb.Result
		link                    gdb.Link
		usedSchema              = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		database, owner, object = tableObjectName(table)
	)
	if database != "" {
		usedSchema = database
	}
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, tableFieldsSqlTmp, owner, object, owner, object)
	if err != nil {
		return nil, err
	}
	if result.IsEmpty() {
		return nil, gerror.NewCodef(
			gcode.CodeNotFound, `table "%s" does not exist in database "%s"`, table, usedSchema,
		)
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		fields[m["Field"].String()] = &gdb.TableField{
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
	return fields, nil
}

// GetTableNameForFields returns the first table in the table expression `table` without its alias,
// keeping its schema and database if it is qualified with them and its quotes if it is quoted, so
// that TableFields looks it up where SQL Server resolves it. It returns an empty string if the first
// table is qualified or quoted but not a table name that splitTableName accepts.
func (d *Driver) GetTableNameForFields(table string) string {
	if tables := gstr.SplitAndTrim(table, ","); len(tables) > 0 {
		first := gstr.SplitAndTrim(tables[0], " ")[0]
		if _, ok := splitTableName(first); ok {
			return first
		}
		if strings.ContainsAny(first, `."[`) {
			return ""
		}
	}
	return d.Core.GetTableNameForFields(table)
}

// tableObjectName returns the database, the schema and the object that the table name `table` is
// made of, as SQL Server stores them without their quote characters, the database and the schema
// being empty if `table` is not qualified with them. `table` is split by splitTableName, or else by
// splitTrimmedTableName, and taken as a single unqualified name if neither accepts it.
func tableObjectName(table string) (database, schema, object string) {
	parts, ok := splitTableName(table)
	if !ok && strings.Contains(table, quoteChar) {
		parts, ok = splitTrimmedTableName(table)
	}
	if !ok {
		return "", "", table
	}
	switch len(parts) {
	case 3:
		return parts[0], parts[1], parts[2]
	case 2:
		return "", parts[0], parts[1]
	default:
		return "", "", parts[0]
	}
}

// splitTrimmedTableName splits `table` into its parts as splitTableName splits it enclosed in double
// quotes, which DriverWrapperDB.TableFields removes from a name whose first and last parts are
// quoted with them, so that `"dbo"."user"` comes as `dbo"."user`. It returns false if `table`
// enclosed in double quotes is not a table name that splitTableName accepts.
func splitTrimmedTableName(table string) (parts []string, ok bool) {
	var (
		quote = quoteChar[0]
		i     int
	)
	for {
		if len(parts) == 0 || (i < len(table) && table[i] == quote) {
			var begin = i
			if len(parts) > 0 {
				begin++
			}
			end := closingQuote(table, begin-1, quote)
			switch {
			case end == len(table)-1:
				return nil, false
			case end < 0:
				end = len(table)
			}
			if end <= begin {
				return nil, false
			}
			parts = append(parts, strings.ReplaceAll(table[begin:end], quoteChar+quoteChar, quoteChar))
			if end == len(table) {
				return parts, true
			}
			i = end + 1
		} else {
			var part string
			if part, i, ok = scanIdentifier(table, i); !ok || i == len(table) {
				return nil, false
			}
			parts = append(parts, part)
		}
		if table[i] != '.' || len(parts) == 3 {
			return nil, false
		}
		i++
	}
}

// splitTableName splits `table`, a table name optionally qualified with its schema, or with its
// database and its schema, like `user`, `dbo.user`, `test.dbo.user`, `"dbo"."user"` or
// `[dbo].[user]`, into its parts as SQL Server stores them, without their quote characters. It
// returns false if `table` is not such a name.
func splitTableName(table string) (parts []string, ok bool) {
	for i := 0; ; i++ {
		var part string
		if part, i, ok = scanIdentifier(table, i); !ok {
			return nil, false
		}
		parts = append(parts, part)
		if i == len(table) {
			return parts, true
		}
		if table[i] != '.' || len(parts) == 3 {
			return nil, false
		}
	}
}

// scanIdentifier returns the identifier that begins at index `i` of `s` without its quote
// characters, being unquoted or quoted with '"' or brackets, and the index following it. It returns
// false if no identifier begins there.
func scanIdentifier(s string, i int) (identifier string, end int, ok bool) {
	if i >= len(s) {
		return "", 0, false
	}
	switch s[i] {
	case '"', '[':
		closing := s[i]
		if closing == '[' {
			closing = ']'
		}
		end = closingQuote(s, i, closing)
		if end <= i+1 {
			return "", 0, false
		}
		quote := string(closing)
		return strings.ReplaceAll(s[i+1:end], quote+quote, quote), end + 1, true
	}
	end = i
	for end < len(s) && isIdentifierChar(s[end]) {
		end++
	}
	return s[i:end], end, end > i
}

// isIdentifierChar reports whether `c` can be part of an unquoted identifier: a letter, a digit,
// `_`, `@`, `#`, `$`, or a byte of a multibyte character.
func isIdentifierChar(c byte) bool {
	return c == '_' || c == '@' || c == '#' || c == '$' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
