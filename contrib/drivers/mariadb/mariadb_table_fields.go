// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	tableFieldsSqlByMariadb = `
SELECT
	c.COLUMN_NAME AS 'Field',
	( CASE WHEN ch.CHECK_CLAUSE LIKE 'json_valid%%' THEN 'json' ELSE c.COLUMN_TYPE END ) AS 'Type',
	c.COLLATION_NAME AS 'Collation',
	c.IS_NULLABLE AS 'Null',
	c.COLUMN_KEY AS 'Key',
	( CASE WHEN c.COLUMN_DEFAULT = 'NULL' OR c.COLUMN_DEFAULT IS NULL THEN NULL ELSE c.COLUMN_DEFAULT END) AS 'Default',
	c.EXTRA AS 'Extra',
	c.PRIVILEGES AS 'Privileges',
	c.COLUMN_COMMENT AS 'Comment' 
FROM
	information_schema.COLUMNS AS c
	LEFT JOIN information_schema.CHECK_CONSTRAINTS AS ch ON c.TABLE_NAME = ch.TABLE_NAME 
	AND c.TABLE_SCHEMA = ch.CONSTRAINT_SCHEMA
	AND c.COLUMN_NAME = ch.CONSTRAINT_NAME 
WHERE
	c.TABLE_SCHEMA = '%s' 
	AND c.TABLE_NAME = '%s'
	ORDER BY c.ORDINAL_POSITION`
)

func init() {
	var err error
	tableFieldsSqlByMariadb, err = gdb.FormatMultiLineSqlToSingle(tableFieldsSqlByMariadb)
	if err != nil {
		panic(err)
	}
}

// TableFields retrieves and returns the fields' information of specified table of current
// schema.
func (d *Driver) TableFields(
	ctx context.Context, table string, schema ...string,
) (fields map[string]*gdb.TableField, err error) {
	var (
		result     gdb.Result
		link       gdb.Link
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}

	result, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(tableFieldsSqlByMariadb, usedSchema, table),
	)
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
