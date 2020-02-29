// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "database/sql"

// batchSqlResult is execution result for batch operations.
type batchSqlResult struct {
	rowsAffected int64
	lastResult   sql.Result
}

// see sql.Result.RowsAffected
func (r *batchSqlResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

// see sql.Result.LastInsertId
func (r *batchSqlResult) LastInsertId() (int64, error) {
	return r.lastResult.LastInsertId()
}
