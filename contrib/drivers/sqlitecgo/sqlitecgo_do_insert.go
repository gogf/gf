// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements SQLiteCGO insert behavior overrides for upsert conflict inference.

package sqlitecgo

import (
	"context"
	"database/sql"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	if option.InsertOption == gdb.InsertOptionSave && len(option.OnConflict) == 0 {
		primaryKeys, err := d.Core.GetPrimaryKeys(ctx, table)
		if err != nil {
			return nil, gerror.WrapCode(
				gcode.CodeInternalError,
				err,
				`failed to get primary keys for Save operation`,
			)
		}
		if !saveDataHasPrimaryKeys(list, primaryKeys) {
			return nil, gerror.NewCodef(
				gcode.CodeMissingParameter,
				`Save operation requires conflict detection: `+
					`either specify OnConflict() columns or ensure table '%s' has primary keys in the data`,
				table,
			)
		}
		option.OnConflict = primaryKeys
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}

// saveDataHasPrimaryKeys reports whether the first save record contains all primary keys.
func saveDataHasPrimaryKeys(list gdb.List, primaryKeys []string) bool {
	if len(list) == 0 || len(primaryKeys) == 0 {
		return false
	}
	for _, primaryKey := range primaryKeys {
		if !saveDataHasKey(list[0], primaryKey) {
			return false
		}
	}
	return true
}

// saveDataHasKey reports whether the save data contains the given key case-insensitively.
func saveDataHasKey(data gdb.Map, key string) bool {
	for dataKey := range data {
		if strings.EqualFold(dataKey, key) {
			return true
		}
	}
	return false
}
