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
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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

	// Check if it is an INSERT statement with primary key.
	if !isUseCoreDoExec && pkField.Name != "" && strings.Contains(strings.ToUpper(sql), "INSERT INTO") {
		primaryKey = pkField.Name
		// Oracle supports RETURNING clause to get the last inserted id
		sql += fmt.Sprintf(returningClause, d.QuoteWord(primaryKey))
	} else {
		// Use default DoExec for non-INSERT or no primary key scenarios
		return d.Core.DoExec(ctx, link, sql, args...)
	}

	// Only the insert operation with primary key can execute the following code

	// SQL filtering.
	sql, args = d.FormatSqlBeforeExecuting(sql, args)
	sql, args, err = d.DoFilter(ctx, link, sql, args)
	if err != nil {
		return nil, err
	}

	// Prepare output variable for RETURNING clause
	var lastInsertId int64
	// Append the output parameter for the RETURNING clause
	args = append(args, &lastInsertId)

	// Link execution.
	_, err = d.DoCommit(ctx, gdb.DoCommitInput{
		Link:          link,
		Sql:           sql,
		Args:          args,
		Stmt:          nil,
		Type:          gdb.SqlTypeExecContext,
		IsTransaction: link.IsTransaction(),
	})

	if err != nil {
		return &Result{
			lastInsertId:      0,
			rowsAffected:      0,
			lastInsertIdError: err,
		}, err
	}

	// Get rows affected from the result
	// For single insert with RETURNING clause, affected is always 1
	var affected int64 = 1

	// Check if the primary key field type supports LastInsertId
	if !strings.Contains(strings.ToLower(pkField.Type), "int") {
		return &Result{
			lastInsertId: 0,
			rowsAffected: affected,
			lastInsertIdError: gerror.NewCodef(
				gcode.CodeNotSupported,
				"LastInsertId is not supported by primary key type: %s",
				pkField.Type,
			),
		}, nil
	}

	return &Result{
		lastInsertId: lastInsertId,
		rowsAffected: affected,
	}, nil
}
