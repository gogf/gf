// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

// Result implements sql.Result interface for Oracle database.
type Result struct {
	lastInsertId      int64
	rowsAffected      int64
	lastInsertIdError error
}

// LastInsertId returns the last insert id.
func (r *Result) LastInsertId() (int64, error) {
	return r.lastInsertId, r.lastInsertIdError
}

// RowsAffected returns the rows affected.
func (r *Result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}
