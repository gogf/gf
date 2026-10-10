// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// DoCommit commits current sql and arguments to underlying sql driver.
// The result of an executed statement is a Result, as the result of the underlying driver
// reports no affected rows.
func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	ctx = d.InjectIgnoreResult(ctx)
	out, err = d.Core.DoCommit(ctx, in)
	if err == nil && (in.Type == gdb.SqlTypeExecContext || in.Type == gdb.SqlTypeStmtExecContext) {
		out.Result = &Result{rowsAffectedErr: errRowsAffectedNotSupported}
	}
	return out, err
}
