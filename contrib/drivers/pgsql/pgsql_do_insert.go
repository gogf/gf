// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// DoInsert inserts or updates data for given table.
// The list parameter must contain at least one record, which was previously validated.
func (d *Driver) DoInsert(
	ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case
		gdb.InsertOptionSave,
		gdb.InsertOptionReplace:
		// PostgreSQL does not support REPLACE INTO syntax, use Save (ON CONFLICT ... DO UPDATE) instead.
		// Automatically detect primary keys if OnConflict is not specified.
		if len(option.OnConflict) == 0 {
			primaryKeys, err := d.Core.GetPrimaryKeys(ctx, table)
			if err != nil {
				return nil, gerror.WrapCode(
					gcode.CodeInternalError,
					err,
					`failed to get primary keys for Save/Replace operation`,
				)
			}
			foundPrimaryKey := false
			for _, conflictKey := range primaryKeys {
				for dataKey := range list[0] {
					if strings.EqualFold(dataKey, conflictKey) {
						foundPrimaryKey = true
						break
					}
				}
				if foundPrimaryKey {
					break
				}
			}
			if !foundPrimaryKey {
				return nil, gerror.NewCodef(
					gcode.CodeMissingParameter,
					`Replace/Save operation requires conflict detection: `+
						`either specify OnConflict() columns or ensure table '%s' has a primary key in the data`,
					table,
				)
			}
			option.OnConflict = primaryKeys
		}
		// Treat Replace as Save operation
		option.InsertOption = gdb.InsertOptionSave

	case gdb.InsertOptionDefault:
		tableFields, err := d.GetCore().GetDB().TableFields(ctx, table)
		if err == nil {
			for _, field := range tableFields {
				if gstr.Equal(field.Key, "pri") {
					pkField := *field
					ctx = context.WithValue(ctx, internalPrimaryKeyInCtx, pkField)
					break
				}
			}
		}

	default:
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}
