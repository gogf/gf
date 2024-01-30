// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql"
)

// InsertIgnore Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) InsertIgnore(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, errUnsupportedInsertIgnore
}

// InsertAndGetId Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) InsertAndGetId(ctx context.Context, table string, data interface{}, batch ...int) (int64, error) {
	return 0, errUnsupportedInsertGetId
}

// Replace Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) Replace(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, errUnsupportedReplace
}
