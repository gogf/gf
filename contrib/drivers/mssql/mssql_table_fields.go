// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
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
INNER JOIN sys.objects o ON c.object_id = o.object_id AND o.type = 'U' AND o.is_ms_shipped = 0
INNER JOIN sys.types t ON c.user_type_id = t.user_type_id
LEFT JOIN sys.default_constraints dc ON c.default_object_id = dc.object_id
LEFT JOIN sys.extended_properties ep ON c.object_id = ep.major_id AND c.column_id = ep.minor_id AND ep.name = 'MS_Description'
LEFT JOIN (
	SELECT ic.object_id, ic.column_id
	FROM sys.index_columns ic
	INNER JOIN sys.indexes i ON ic.object_id = i.object_id AND ic.index_id = i.index_id
	WHERE i.is_primary_key = 1
) pk ON c.object_id = pk.object_id AND c.column_id = pk.column_id
WHERE o.name = '%s'
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
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result     gdb.Result
		link       gdb.Link
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	structureSql := fmt.Sprintf(tableFieldsSqlTmp, table)
	result, err = d.DoSelect(ctx, link, structureSql)
	if err != nil {
		return nil, err
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
