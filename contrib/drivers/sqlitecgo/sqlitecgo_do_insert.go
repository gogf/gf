// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// sqlitecgo_do_insert.go infers Save conflict columns from table primary keys.

package sqlitecgo

import (
	"context"
	"database/sql"

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
		if !gdb.HasPrimaryKeys(list, primaryKeys) {
			return nil, gerror.NewCodef(
				gcode.CodeMissingParameter,
				`Save operation requires conflict detection: `+
					`either specify OnConflict() columns or include all primary key values for table '%s' in the save data`,
				table,
			)
		}
		option.OnConflict = primaryKeys
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}
