// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
)

// DoCommit commits current sql and arguments to underlying sql driver.
// It executes an insert that doInsertDefault marks with the RETURNING clause of the identity
// column, and returns the value of the column as the last insert id.
func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	if in.Type == gdb.SqlTypeExecContext {
		if column, ok := ctx.Value(ctxKeyReturningColumn).(string); ok {
			return d.doExecReturning(ctx, in, column)
		}
	}
	return d.Core.DoCommit(ctx, in)
}

// doExecReturning executes the insert of `in` with the clause returning `column` into an output
// parameter, which becomes the last insert id of the result.
func (d *Driver) doExecReturning(ctx context.Context, in gdb.DoCommitInput, column string) (out gdb.DoCommitOutput, err error) {
	var id int64
	in.Sql += fmt.Sprintf(" RETURNING %s INTO ?", d.QuoteWord(column))
	in.Link = &returningLink{Link: in.Link, dest: &id}
	if out, err = d.Core.DoCommit(ctx, in); err != nil {
		return out, err
	}
	out.Result = &Result{Result: out.Result, lastInsertId: id}
	return out, nil
}

// returningLink is a gdb.Link that appends the output parameter of the RETURNING clause
// to the arguments of the statement it executes.
type returningLink struct {
	gdb.Link
	dest *int64
}

// ExecContext executes the statement with the RETURNING output parameter appended to `args`.
func (l *returningLink) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return l.Link.ExecContext(ctx, query, append(args, sql.Out{Dest: l.dest})...)
}
