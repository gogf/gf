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
	"strings"
)

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by pgsql driver`,
		)

	case gdb.InsertOptionIgnore:
		return d.doInsertIgnore(ctx, link, table, list)

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

// doInsertIgnore inserts a list of records into the specified table in pgsql.
// INSERT INTO <table> (<columns>) VALUES (<values>) ON CONFLICT DO NOTHING
func (d *Driver) doInsertIgnore(ctx context.Context, link gdb.Link, table string, list gdb.List) (result sql.Result, err error) {
	if len(list) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidRequest, `Insert operation list is empty for pgsql driver`)
	}

	var (
		one          = list[0]
		charL, charR = d.GetChars()
		insertKeys   = make([]string, 0, len(one))
		insertValues = make([]string, 0, len(one))
		queryValues  = make([]interface{}, 0, len(one))
	)

	for key, value := range one {
		insertKeys = append(insertKeys, charL+key+charR)
		insertValues = append(insertValues, "?")
		queryValues = append(queryValues, value)
	}

	sqlStr := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING`,
		table,
		strings.Join(insertKeys, ","),
		strings.Join(insertValues, ","),
	)

	result, err = d.DoExec(ctx, link, sqlStr, queryValues...)
	if err != nil {
		return result, err
	}

	return result, nil
}
