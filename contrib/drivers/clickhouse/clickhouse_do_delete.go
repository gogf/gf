// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql"

	"github.com/gogf/gf/v2/database/gdb"
)

// DoDelete does "DELETE FROM ... " statement for the table.
func (d *Driver) DoDelete(ctx context.Context, link gdb.Link, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	ctx = d.injectNeedParsedSql(ctx)
	return d.Core.DoDelete(ctx, link, table, condition, args...)
}
