// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

// GetBoolLiteral returns the SQL literal for the given boolean value.
// ClickHouse has a strict Bool type that requires 'true'/'false' literals.
// GetLockSharedClause is not overridden because ClickHouse does not support
// row-level locking — queries with lock clauses would fail at the server
// regardless of which syntax is emitted.
func (d *Driver) GetBoolLiteral(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
