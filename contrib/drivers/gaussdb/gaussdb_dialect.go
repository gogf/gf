// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// GetBoolLiteral returns the SQL literal for the given boolean value.
// GaussDB is PostgreSQL-compatible and uses strict boolean literals.
func (d *Driver) GetBoolLiteral(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// GetLockSharedClause returns the SQL clause for shared row locks.
// GaussDB uses "FOR SHARE" instead of MySQL's legacy "LOCK IN SHARE MODE".
func (d *Driver) GetLockSharedClause() string {
	return gdb.LockForShare
}
