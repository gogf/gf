// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"github.com/gogf/gf/v2/database/gdb"
)

// GetBoolLiteral returns the SQL literal for the given boolean value.
// PostgreSQL uses a strict boolean type that requires 'true'/'false' literals,
// not MySQL's numeric 1/0.
func (d *Driver) GetBoolLiteral(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// GetLockSharedClause returns the SQL clause for shared row locks.
// PostgreSQL uses "FOR SHARE" instead of MySQL's legacy "LOCK IN SHARE MODE".
func (d *Driver) GetLockSharedClause() string {
	return gdb.LockForShare
}

// GetPartitionClause returns an empty string for PostgreSQL. PostgreSQL does
// not support the MySQL-style PARTITION (p) clause on SELECT statements;
// it uses declarative partitioning with inheritance-based partition pruning.
func (d *Driver) GetPartitionClause(partitions string) string {
	return ""
}
