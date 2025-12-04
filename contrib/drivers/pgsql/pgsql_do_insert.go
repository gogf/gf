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
	"github.com/gogf/gf/v2/text/gstr"
)

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionReplace:
		// PostgreSQL does not support REPLACE INTO syntax, use Save (ON CONFLICT ... DO UPDATE) instead.
		// Automatically detect primary keys if OnConflict is not specified.
		if len(option.OnConflict) == 0 {
			primaryKeys, err := d.getPrimaryKeys(ctx, table)
			if err != nil {
				return nil, gerror.WrapCode(
					gcode.CodeInternalError,
					err,
					`failed to get primary keys for Replace operation`,
				)
			}
			if len(primaryKeys) == 0 {
				return nil, gerror.NewCode(
					gcode.CodeMissingParameter,
					`Replace operation requires primary key, but table has no primary key defined`,
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
				if field.Key == "pri" {
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

// getPrimaryKeys retrieves the primary key field list of the table.
// This method extracts primary key information from TableFields.
func (d *Driver) getPrimaryKeys(ctx context.Context, table string) ([]string, error) {
	tableFields, err := d.TableFields(ctx, table)
	if err != nil {
		return nil, err
	}

	var primaryKeys []string
	for _, field := range tableFields {
		if gstr.Equal(field.Key, "PRI") {
			primaryKeys = append(primaryKeys, field.Name)
		}
	}

	return primaryKeys, nil
}
