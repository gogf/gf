// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// DoInsert inserts or updates data forF given table.
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by pgsql driver`,
		)

	case gdb.InsertOptionIgnore:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Insert ignore operation is not supported by pgsql driver`,
		)

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
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}

// FormatUpsert returns SQL clause of type upsert for PgSQL.
// For example: ON CONFLICT (id) DO UPDATE SET ...
func (d *Driver) FormatUpsert(columns []string, list gdb.List, option gdb.DoInsertOption) (string, error) {
	if len(option.OnConflict) == 0 {
		return "", gerror.New("Please specify conflict columns")
	}

	var onDuplicateStr string
	if option.OnDuplicateStr != "" {
		onDuplicateStr = option.OnDuplicateStr
	} else if len(option.OnDuplicateMap) > 0 {
		for k, v := range option.OnDuplicateMap {
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			switch v.(type) {
			case gdb.Raw, *gdb.Raw:
				onDuplicateStr += fmt.Sprintf(
					"%s=%s",
					d.Core.QuoteWord(k),
					v,
				)
			default:
				onDuplicateStr += fmt.Sprintf(
					"%s=EXCLUDED.%s",
					d.Core.QuoteWord(k),
					d.Core.QuoteWord(gconv.String(v)),
				)
			}
		}
	} else {
		for _, column := range columns {
			// If it's SAVE operation, do not automatically update the creating time.
			if d.Core.IsSoftCreatedFieldName(column) {
				continue
			}
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			onDuplicateStr += fmt.Sprintf(
				"%s=EXCLUDED.%s",
				d.Core.QuoteWord(column),
				d.Core.QuoteWord(column),
			)
		}
	}

	conflictKeys := gstr.Join(option.OnConflict, ",")

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET ", conflictKeys) + onDuplicateStr, nil
}
