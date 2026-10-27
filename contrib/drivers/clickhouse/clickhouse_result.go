// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	errLastInsertIdNotSupported = gerror.NewCode(
		gcode.CodeNotSupported, "LastInsertId is not supported by ClickHouse",
	)
	errRowsAffectedNotSupported = gerror.NewCode(
		gcode.CodeNotSupported, "RowsAffected is not supported by ClickHouse for this statement",
	)
)

// Result is the sql.Result of a statement the driver executes. ClickHouse reports neither the
// id of an inserted row nor the rows that a mutation affects, so that LastInsertId always
// fails as not supported, and RowsAffected fails unless the driver knows the rows itself, as
// it does for the records it inserts.
type Result struct {
	rowsAffected    int64
	rowsAffectedErr error
}

// LastInsertId returns an error of code gcode.CodeNotSupported.
func (r *Result) LastInsertId() (int64, error) {
	return 0, errLastInsertIdNotSupported
}

// RowsAffected returns the number of records inserted, or an error of code
// gcode.CodeNotSupported for any other statement.
func (r *Result) RowsAffected() (int64, error) {
	return r.rowsAffected, r.rowsAffectedErr
}
