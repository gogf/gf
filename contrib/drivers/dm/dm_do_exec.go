// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"regexp"

	"github.com/gogf/gf/v2/database/gdb"
)

// releaseSavePointRegex matches the statement releasing a savepoint, which the core issues when a
// nested transaction commits.
var releaseSavePointRegex = regexp.MustCompile(`(?i)^\s*RELEASE\s+SAVEPOINT\s+"?\w+"?\s*$`)

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
// It does not execute the statement releasing a savepoint, which DM does not support: DM keeps a
// savepoint until the transaction ends, so that a nested transaction committing leaves it to the
// outer transaction as releasing it would.
func (d *Driver) DoExec(ctx context.Context, link gdb.Link, sql string, args ...any) (result sql.Result, err error) {
	if releaseSavePointRegex.MatchString(sql) {
		return driver.RowsAffected(0), nil
	}
	return d.Core.DoExec(ctx, link, sql, args...)
}
