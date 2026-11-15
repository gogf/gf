// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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
			if !gdb.HasPrimaryKeys(list, primaryKeys) {
				return nil, gerror.NewCodef(
					gcode.CodeMissingParameter,
					`Replace/Save operation requires conflict detection: `+
						`either specify OnConflict() columns or include all primary key values for table '%s' in the save data`,
					table,
				)
			}
			option.OnConflict = primaryKeys
		}
		// Treat Replace as Save operation
		option.InsertOption = gdb.InsertOptionSave

	// pgsql support InsertIgnore natively, so no need to set primary key in context.
	case gdb.InsertOptionIgnore, gdb.InsertOptionDefault:
		// DoExec appends RETURNING for this field and uses its type to decide whether
		// LastInsertId is supported. The column is the generated primary key, which is
		// not always the first key column in table order.
		tableFields, err := d.GetCore().GetDB().TableFields(ctx, table)
		if err == nil {
			if field := gdb.InsertIDField(gdb.PrimaryKeyFields(tableFields)); field != nil {
				ctx = context.WithValue(ctx, internalPrimaryKeyInCtx, *field)
			}
		}

	default:
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}
