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

// DoUpdate does "UPDATE ... " statement for the table.
// The values of the Decimal columns in a map `data` are converted as convertDecimalValues does.
func (d *Driver) DoUpdate(ctx context.Context, link gdb.Link, table string, data any, condition string, args ...any) (result sql.Result, err error) {
	if record, ok := data.(gdb.Map); ok {
		fields, err := d.tableFields(ctx, table)
		if err != nil {
			return nil, err
		}
		records, err := d.convertDecimalValues(ctx, fields, gdb.List{record})
		if err != nil {
			return nil, err
		}
		data = records[0]
	}
	ctx = d.injectNeedParsedSql(ctx)
	return d.Core.DoUpdate(ctx, link, table, data, condition, args...)
}
