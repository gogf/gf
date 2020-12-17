// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "database/sql"

// SqlResult is execution result for sql operations.
// It also supports batch operation result for rowsAffected.
type SqlResult struct {
	result   sql.Result
	affected int64
}

// MustGetAffected returns the affected rows count, if any error occurs, it panics.
func (r *SqlResult) MustGetAffected() int64 {
	rows, err := r.RowsAffected()
	if err != nil {
		panic(err)
	}
	return rows
}

// MustGetInsertId returns the last insert id, if any error occurs, it panics.
func (r *SqlResult) MustGetInsertId() int64 {
	id, err := r.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}

// see sql.Result.RowsAffected
func (r *SqlResult) RowsAffected() (int64, error) {
	if r.affected > 0 {
		return r.affected, nil
	}
	if r.result == nil {
		return 0, nil
	}
	return r.result.RowsAffected()
}

// see sql.Result.LastInsertId
func (r *SqlResult) LastInsertId() (int64, error) {
	if r.result == nil {
		return 0, nil
	}
	return r.result.LastInsertId()
}
