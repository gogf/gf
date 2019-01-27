// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import "database/sql"

// 批量执行的结果对象
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