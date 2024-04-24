// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/intlog"
)

// RowsToResult converts underlying data record type sql.Rows to Result type.
func (c *Core) RowsToResult(ctx context.Context, rows *sql.Rows) (Result, error) {
	if rows == nil {
		return nil, nil
	}
	defer func() {
		if err := rows.Close(); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}()
	if !rows.Next() {
		return nil, nil
	}
	// Column names and types.
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	if len(columnTypes) > 0 {
		if internalData := c.getInternalColumnFromCtx(ctx); internalData != nil {
			internalData.FirstResultColumn = columnTypes[0].Name()
		}
	}
	var (
		values   = make([]interface{}, len(columnTypes))
		result   = make(Result, 0)
		scanArgs = make([]interface{}, len(values))
	)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for {
		if err = rows.Scan(scanArgs...); err != nil {
			return result, err
		}
		record := Record{}
		for i, value := range values {
			if value == nil {
				// DO NOT use `gvar.New(nil)` here as it creates an initialized object
				// which will cause struct converting issue.
				record[columnTypes[i].Name()] = nil
			} else {
				var convertedValue interface{}
				if convertedValue, err = c.columnValueToLocalValue(ctx, value, columnTypes[i]); err != nil {
					return nil, err
				}
				record[columnTypes[i].Name()] = gvar.New(convertedValue)
			}
		}
		result = append(result, record)
		if !rows.Next() {
			break
		}
	}
	return result, nil
}
