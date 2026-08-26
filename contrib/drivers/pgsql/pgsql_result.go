// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file defines PostgreSQL execution results and generated primary-key values.

package pgsql

import (
	"database/sql"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
)

// Result implements SQL execution result access for PostgreSQL statements.
type Result struct {
	sql.Result
	affected                  int64     // affected stores the number of returned rows.
	lastInsertId              int64     // lastInsertId stores an integer primary key.
	lastInsertIdError         error     // lastInsertIdError reports unsupported primary-key types.
	lastInsertPrimaryKeyValue gdb.Value // lastInsertPrimaryKeyValue stores the returned primary key.
}

// RowsAffected returns the number of rows returned by the executed statement.
func (pgr Result) RowsAffected() (int64, error) {
	return pgr.affected, nil
}

// LastInsertId returns an integer primary key or the recorded unsupported error.
func (pgr Result) LastInsertId() (int64, error) {
	return pgr.lastInsertId, pgr.lastInsertIdError
}

// LastInsertPrimaryKeyValue returns the primary key value produced by INSERT ... RETURNING.
func (pgr Result) LastInsertPrimaryKeyValue() (gdb.Value, error) {
	if pgr.lastInsertPrimaryKeyValue == nil {
		return gvar.New(nil), nil
	}
	return pgr.lastInsertPrimaryKeyValue, nil
}
