// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import "database/sql"

type Result struct {
	sql.Result
	affected          int64
	lastInsertId      int64
	lastInsertIdError error
}

func (pgr Result) RowsAffected() (int64, error) {
	return pgr.affected, nil
}

func (pgr Result) LastInsertId() (int64, error) {
	return pgr.lastInsertId, pgr.lastInsertIdError
}
