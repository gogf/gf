// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

const (
	returningClause = " RETURNING %s INTO ?"
)

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
// It handles INSERT statements specially to support LastInsertId.
func (d *Driver) DoExec(
	ctx context.Context, link gdb.Link, sql string, args ...interface{},
) (result sql.Result, err error) {
	var (
		isUseCoreDoExec = true
		primaryKey      string
		pkField         gdb.TableField
	)

	// Transaction checks.
	if link == nil {
		if tx := gdb.TXFromCtx(ctx, d.GetGroup()); tx != nil {
			link = tx
		} else if link, err = d.MasterLink(); err != nil {
			return nil, err
		}
	} else if !link.IsTransaction() {
		if tx := gdb.TXFromCtx(ctx, d.GetGroup()); tx != nil {
			link = tx
		}
	}

	// Check if it is an insert operation with primary key from context.
	if value := ctx.Value(internalPrimaryKeyInCtx); value != nil {
		if field, ok := value.(gdb.TableField); ok {
			pkField = field
			isUseCoreDoExec = false
		}
	}

	// Check if it is an INSERT statement with an integer primary key.
	if !isUseCoreDoExec && pkField.Name != "" && isIntegerField(pkField) &&
		strings.Contains(strings.ToUpper(sql), "INSERT INTO") {
		primaryKey = pkField.Name
		// Oracle supports RETURNING clause to get the last inserted id
		sql += fmt.Sprintf(returningClause, d.QuoteWord(primaryKey))
	} else {
		// Use default DoExec for non-INSERT, no primary key or non-integer primary key scenarios
		return d.Core.DoExec(ctx, link, sql, args...)
	}

	// Only the insert operation with primary key can execute the following code

	// Prepare output variable for RETURNING clause
	var lastInsertId int64
	r, err := d.Core.DoExec(ctx, &returningLink{Link: link, dest: &lastInsertId}, sql, args...)
	if err != nil {
		return &Result{
			lastInsertId:      0,
			rowsAffected:      0,
			lastInsertIdError: err,
		}, err
	}
	affected, err := r.RowsAffected()
	if err != nil {
		return nil, err
	}
	return &Result{
		lastInsertId: lastInsertId,
		rowsAffected: affected,
	}, nil
}

// isIntegerField reports whether the column `field` holds integers, which TableFields reports
// as INT(precision,scale) for a NUMBER column without fractional digits.
func isIntegerField(field gdb.TableField) bool {
	typeName, _, _ := strings.Cut(field.Type, "(")
	return strings.EqualFold(strings.TrimSpace(typeName), "INT")
}

// returningLink is a gdb.Link that appends the output parameter of the RETURNING clause
// to the arguments of the statement it executes.
type returningLink struct {
	gdb.Link
	dest *int64
}

// ExecContext executes the statement with the RETURNING output parameter appended to `args`.
func (l *returningLink) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return l.Link.ExecContext(ctx, query, append(args, l.dest)...)
}
