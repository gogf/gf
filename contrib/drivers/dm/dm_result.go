// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"database/sql"
)

// Result is the sql.Result of an insert, whose last insert id is the value of the identity
// column that the driver determines rather than the one the underlying driver reports.
type Result struct {
	sql.Result
	lastInsertId int64
}

// LastInsertId returns the last insert id.
func (r *Result) LastInsertId() (int64, error) {
	return r.lastInsertId, nil
}
