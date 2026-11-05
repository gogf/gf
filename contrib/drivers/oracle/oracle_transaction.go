// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"regexp"
)

// releaseSavePointRegex matches the statement releasing a savepoint, which the core issues when a
// nested transaction commits.
var releaseSavePointRegex = regexp.MustCompile(`(?i)^\s*RELEASE\s+SAVEPOINT\s+"?\w+"?\s*$`)

// isReleaseSavePoint reports whether `sql` releases a savepoint. Oracle has no such statement and
// keeps a savepoint until the outer transaction ends, so that the statement is not executed.
func isReleaseSavePoint(sql string) bool {
	return releaseSavePointRegex.MatchString(sql)
}
