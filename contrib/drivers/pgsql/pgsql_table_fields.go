// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	tableFieldsSqlTmp = `
SELECT
    a.attname                                                                            AS field,
    t.typname                                                                            AS type,
    a.attnotnull                                                                         AS null,
    (CASE WHEN d.contype = 'p' THEN 'pri' WHEN d.contype = 'u' THEN 'uni' ELSE '' END)   AS key,
    ic.column_default                                                                    AS default_value,
    b.description                                                                        AS comment,
    COALESCE(character_maximum_length, numeric_precision, -1)                            AS length,
    numeric_scale                                                                        AS scale
FROM pg_attribute a
    LEFT JOIN pg_class c                 ON a.attrelid = c.oid
    LEFT JOIN pg_constraint d            ON d.conrelid = c.oid AND a.attnum = d.conkey[1]
    LEFT JOIN pg_description b           ON a.attrelid = b.objoid AND a.attnum = b.objsubid
    LEFT JOIN pg_type t                  ON a.atttypid = t.oid
    LEFT JOIN information_schema.columns ic ON ic.column_name = a.attname AND ic.table_name = c.relname
WHERE c.oid = '%s'::regclass
    AND a.attisdropped IS FALSE
    AND a.attnum > 0
ORDER BY a.attnum`
)

func init() {
	var err error
	tableFieldsSqlTmp, err = gdb.FormatMultiLineSqlToSingle(tableFieldsSqlTmp)
	if err != nil {
		panic(err)
	}
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result     gdb.Result
		link       gdb.Link
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		// TODO duplicated `id` result?
		structureSql = fmt.Sprintf(tableFieldsSqlTmp, table)
	)
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, structureSql)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	var (
		index         = 0
		name          string
		ok            bool
		existingField *gdb.TableField
	)
	for _, m := range result {
		name = m["field"].String()
		// Merge duplicated fields, especially for key constraints.
		// Priority: pri > uni > others
		if existingField, ok = fields[name]; ok {
			currentKey := m["key"].String()
			// Merge key information with priority: pri > uni
			if currentKey == "pri" || (currentKey == "uni" && existingField.Key != "pri") {
				existingField.Key = currentKey
			}
			continue
		}
		fields[name] = &gdb.TableField{
			Index:   index,
			Name:    name,
			Type:    m["type"].String(),
			Null:    !m["null"].Bool(),
			Key:     m["key"].String(),
			Default: m["default_value"].Val(),
			Comment: m["comment"].String(),
		}
		index++
	}
	return fields, nil
}
