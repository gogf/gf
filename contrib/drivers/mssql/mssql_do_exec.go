// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// outputInsertedClause is the OUTPUT clause returning a column of the inserted rows.
const outputInsertedClause = " OUTPUT INSERTED.%s"

// The constants of GetInsertOutputSql.
//
// Deprecated: GetInsertOutputSql is deprecated.
const (
	// Database field attributes
	fieldKeyPrimary = "PRI"

	// SQL keywords and syntax markers
	outputKeyword = "OUTPUT"

	// Object and field references
	insertedObjectName = "INSERTED"

	// Result field names and aliases
	affectCountExpression  = " 1 as AffectCount"
	lastInsertIdFieldAlias = "ID"
)

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
// The INSERT statements of DoInsert return the IDENTITY column of the inserted rows through an
// OUTPUT clause, whose first value is reported as LastInsertId.
func (d *Driver) DoExec(ctx context.Context, link gdb.Link, sqlStr string, args ...any) (result sql.Result, err error) {
	if isReleaseSavePoint(sqlStr) {
		return driver.RowsAffected(0), nil
	}
	identityField, isInsert := ctx.Value(internalIdentityFieldInCtx).(string)
	if !isInsert {
		return d.Core.DoExec(ctx, link, sqlStr, args...)
	}
	pos := insertColumnListEnd(sqlStr)
	if identityField == "" || pos < 0 {
		r, err := d.Core.DoExec(ctx, link, sqlStr, args...)
		if err != nil {
			return r, err
		}
		affected, err := r.RowsAffected()
		if err != nil {
			return nil, err
		}
		return &Result{rowsAffected: affected}, nil
	}

	// Transaction checks.
	if link == nil {
		if tx := gdb.TXFromCtx(ctx, d.GetGroup()); tx != nil {
			// Firstly, check and retrieve transaction link from context.
			link = &txLinkMssql{tx.GetSqlTX()}
		} else if link, err = d.MasterLink(); err != nil {
			// Or else it creates one from master node.
			return nil, err
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
		if tx := gdb.TXFromCtx(ctx, d.GetGroup()); tx != nil {
			link = &txLinkMssql{tx.GetSqlTX()}
		}
	}
	sqlStr = sqlStr[:pos] + fmt.Sprintf(outputInsertedClause, d.QuoteWord(identityField)) + sqlStr[pos:]
	return d.Core.DoExec(ctx, &outputLink{Link: link}, sqlStr, args...)
}

// insertColumnListEnd returns the position following the parenthesis that closes the column list
// of the INSERT statement `sqlStr` built by gdb.Core.DoInsert, or -1 if there is none.
// Parentheses inside quoted identifiers are skipped.
func insertColumnListEnd(sqlStr string) int {
	var (
		depth  int
		quoted bool
	)
	for i := 0; i < len(sqlStr); i++ {
		if sqlStr[i] == quoteChar[0] {
			quoted = !quoted
			continue
		}
		if quoted {
			continue
		}
		switch sqlStr[i] {
		case '(':
			depth++
		case ')':
			if depth--; depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

// GetTableNameFromSql get table name from sql statement
// It handles table string like:
// "user"
// "user u"
// "DbLog.dbo.user",
// "user as u".
//
// Deprecated: DoExec no longer uses it, and it will be removed in the future.
func (d *Driver) GetTableNameFromSql(sqlStr string) (table string) {
	// INSERT INTO "ip_to_id"("ip") OUTPUT  1 as AffectCount,INSERTED.id as ID VALUES(?)
	var (
		leftChars, rightChars = d.GetChars()
		trimStr               = leftChars + rightChars + "[] "
		pattern               = "INTO(.+?)\\("
		regCompile            = regexp.MustCompile(pattern)
		tableInfo             = regCompile.FindStringSubmatch(sqlStr)
	)
	// get the first one. after the first it may be content of the value, it's not table name.
	table = tableInfo[1]
	table = strings.Trim(table, " ")
	if strings.Contains(table, ".") {
		tmpAry := strings.Split(table, ".")
		// the last one is table name
		table = tmpAry[len(tmpAry)-1]
	} else if strings.Contains(table, "as") || strings.Contains(table, " ") {
		tmpAry := strings.Split(table, "as")
		if len(tmpAry) < 2 {
			tmpAry = strings.Split(table, " ")
		}
		// get the first one
		table = tmpAry[0]
	}
	table = strings.Trim(table, trimStr)
	return table
}

// txLink is used to implement interface Link for TX.
type txLinkMssql struct {
	*sql.Tx
}

// IsTransaction returns if current Link is a transaction.
func (l *txLinkMssql) IsTransaction() bool {
	return true
}

// IsOnMaster checks and returns whether current link is operated on master node.
// Note that, transaction operation is always operated on master node.
func (l *txLinkMssql) IsOnMaster() bool {
	return true
}

// GetInsertOutputSql  gen get last_insert_id code
//
// Deprecated: DoExec no longer uses it, and it will be removed in the future.
func (d *Driver) GetInsertOutputSql(ctx context.Context, table string) string {
	fds, errFd := d.GetDB().TableFields(ctx, table)
	if errFd != nil {
		return ""
	}
	extraSqlAry := make([]string, 0)
	extraSqlAry = append(extraSqlAry, fmt.Sprintf(" %s %s", outputKeyword, affectCountExpression))
	incrNo := 0
	if len(fds) > 0 {
		for _, fd := range fds {
			// has primary key and is auto-increment
			if fd.Extra == fieldExtraIdentity && fd.Key == fieldKeyPrimary && !fd.Null {
				incrNoStr := ""
				if incrNo == 0 { // fixed first field named id, convenient to get
					incrNoStr = fmt.Sprintf(" as %s", lastInsertIdFieldAlias)
				}

				extraSqlAry = append(extraSqlAry, fmt.Sprintf("%s.%s%s", insertedObjectName, fd.Name, incrNoStr))
				incrNo++
			}
			// fmt.Printf("null:%t name:%s key:%s k:%s \n", fd.Null, fd.Name, fd.Key, k)
		}
	}
	return strings.Join(extraSqlAry, ",")
	// sql example:INSERT INTO "ip_to_id"("ip") OUTPUT  1 as AffectCount,INSERTED.id as ID VALUES(?)
}

// outputLink is a gdb.Link that executes an INSERT statement carrying an OUTPUT clause as a query,
// as SQL Server returns the OUTPUT values as a result set.
type outputLink struct {
	gdb.Link
}

// ExecContext executes the INSERT statement `query` and returns a result reporting the count of
// the rows it returns as RowsAffected and the value of the first one as LastInsertId.
func (l *outputLink) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	rows, err := l.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := new(Result)
	for rows.Next() {
		if result.rowsAffected == 0 {
			if err = rows.Scan(&result.lastInsertId); err != nil {
				return nil, err
			}
		}
		result.rowsAffected++
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
