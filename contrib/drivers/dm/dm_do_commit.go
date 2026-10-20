// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// isolationLevelNames maps the isolation levels that DM supports to their names in the
// SET TRANSACTION statement.
var isolationLevelNames = map[sql.IsolationLevel]string{
	sql.LevelReadUncommitted: "READ UNCOMMITTED",
	sql.LevelReadCommitted:   "READ COMMITTED",
	sql.LevelRepeatableRead:  "REPEATABLE READ",
	sql.LevelSerializable:    "SERIALIZABLE",
}

// DoCommit commits current sql and arguments to underlying sql driver.
//
// It begins a transaction with the isolation level and the access mode of its options set by
// SET TRANSACTION statements, which last for the transaction only, as they do on MySQL. The
// underlying driver sets them on the session instead, which keeps them on the pooled connection
// after the transaction ends. An isolation level DM does not support is still passed to the
// underlying driver, which rejects it.
//
// It executes an insert that doInsertReturning marks with the RETURNING clause of the primary key,
// and returns the value of the key as the last insert id.
func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	switch in.Type {
	case gdb.SqlTypeBegin:
		return d.doBegin(ctx, in)

	case gdb.SqlTypeExecContext:
		if column, ok := ctx.Value(ctxKeyReturningColumn).(string); ok {
			return d.doExecReturning(ctx, in, column)
		}
	default:
	}
	return d.Core.DoCommit(ctx, in)
}

// doBegin begins a transaction with the options of `in`, see DoCommit. It sets both options by one
// SET TRANSACTION statement, as a second one would apply to the transaction after this one.
func (d *Driver) doBegin(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	var options []string
	if name, ok := isolationLevelNames[in.TxOptions.Isolation]; ok {
		options = append(options, "ISOLATION LEVEL "+name)
		in.TxOptions.Isolation = sql.LevelDefault
	}
	if in.TxOptions.ReadOnly {
		options = append(options, "READ ONLY")
		in.TxOptions.ReadOnly = false
	}
	if out, err = d.Core.DoCommit(ctx, in); err != nil || len(options) == 0 {
		return out, err
	}
	// The options apply to the transaction that the next statement begins, so a statement begins
	// it, as a transaction executing none would leave them to the next transaction.
	for _, statement := range []string{"SET TRANSACTION " + strings.Join(options, ", "), "SELECT 1 FROM DUAL"} {
		if _, err = out.Tx.Exec(statement); err != nil {
			if rollbackErr := out.Tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
			return gdb.DoCommitOutput{}, err
		}
	}
	return out, nil
}

// doExecReturning executes the insert of `in` with the clause returning `column` into an output
// parameter, which becomes the last insert id of the result. The parameter is a sql.NullInt64, as
// the underlying driver fails to convert a NUMBER column into an *int64 (error 6007).
func (d *Driver) doExecReturning(ctx context.Context, in gdb.DoCommitInput, column string) (out gdb.DoCommitOutput, err error) {
	var id sql.NullInt64
	in.Sql = strings.TrimRight(in.Sql, " ") + fmt.Sprintf(" RETURNING %s INTO ?", d.QuoteWord(column))
	in.Link = &returningLink{Link: in.Link, dest: &id}
	if out, err = d.Core.DoCommit(ctx, in); err != nil {
		return out, err
	}
	out.Result = &insertResult{Result: out.Result, lastInsertId: id.Int64}
	return out, nil
}

// returningLink is a gdb.Link that appends the output parameter of the RETURNING clause
// to the arguments of the statement it executes.
type returningLink struct {
	gdb.Link
	dest *sql.NullInt64
}

// ExecContext executes the statement with the RETURNING output parameter appended to `args`.
func (l *returningLink) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return l.Link.ExecContext(ctx, query, append(args, sql.Out{Dest: l.dest})...)
}
