// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"regexp"
)

var (
	// savePointRegex matches the statement creating a savepoint, which the core issues for a
	// nested transaction and for TX.SavePoint.
	savePointRegex = regexp.MustCompile(`(?i)^\s*SAVEPOINT\s+(\S+)\s*$`)
	// rollbackToSavePointRegex matches the statement rolling back to a savepoint, which the core
	// issues when a nested transaction rolls back and for TX.RollbackTo.
	rollbackToSavePointRegex = regexp.MustCompile(`(?i)^\s*ROLLBACK\s+TO\s+SAVEPOINT\s+(\S+)\s*$`)
	// releaseSavePointRegex matches the statement releasing a savepoint, which the core issues
	// when a nested transaction commits.
	releaseSavePointRegex = regexp.MustCompile(`(?i)^\s*RELEASE\s+SAVEPOINT\s+\S+\s*$`)
)

// rewriteSavePoint returns the statement `sql` creating or rolling back to a savepoint in SQL
// Server's syntax, SAVE TRANSACTION or ROLLBACK TRANSACTION, or `sql` itself if it is neither.
func rewriteSavePoint(sql string) string {
	if match := savePointRegex.FindStringSubmatch(sql); match != nil {
		return "SAVE TRANSACTION " + match[1]
	}
	if match := rollbackToSavePointRegex.FindStringSubmatch(sql); match != nil {
		return "ROLLBACK TRANSACTION " + match[1]
	}
	return sql
}

// isReleaseSavePoint reports whether `sql` releases a savepoint. SQL Server has no such
// statement and keeps a savepoint until the outer transaction ends, so that it is not executed.
func isReleaseSavePoint(sql string) bool {
	return releaseSavePointRegex.MatchString(sql)
}
