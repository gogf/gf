// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	// Special insert/ignore operation for sqlite.
	switch {
	case gstr.HasPrefix(sql, gdb.InsertOperationIgnore):
		sql = "INSERT OR IGNORE" + sql[len(gdb.InsertOperationIgnore):]

	case gstr.HasPrefix(sql, gdb.InsertOperationReplace):
		sql = "INSERT OR REPLACE" + sql[len(gdb.InsertOperationReplace):]

	default:
		if gstr.Contains(sql, gdb.InsertOnDuplicateKeyUpdate) {
			return sql, args, gerror.NewCode(
				gcode.CodeNotSupported,
				`Save operation is not supported by sqlite driver`,
			)
		}
	}
	return d.Core.DoFilter(ctx, link, sql, args)
}
