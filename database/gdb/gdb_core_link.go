// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
)

// dbLink is used to implement interface Link for DB.
type dbLink struct {
	*sql.DB
}

// txLink is used to implement interface Link for TX.
type txLink struct {
	*sql.Tx
}

// IsTransaction returns if current Link is a transaction.
func (*dbLink) IsTransaction() bool {
	return false
}

// IsTransaction returns if current Link is a transaction.
func (*txLink) IsTransaction() bool {
	return true
}
