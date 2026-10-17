// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"strings"
)

// FoldIdentifier returns `name` in upper case, as SQL Server matches identifiers case-insensitively.
func (d *Driver) FoldIdentifier(name string) string {
	return strings.ToUpper(name)
}
