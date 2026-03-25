// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

// Result instance of sql.Result
type Result struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (r *Result) LastInsertId() (int64, error) {
	return r.lastInsertId, r.err
}

func (r *Result) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}
