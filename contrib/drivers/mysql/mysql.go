// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package mysql implements gdb.Driver, which supports operations for MySQL.
package mysql

import (
	"github.com/gogf/gf/v2/database/gdb"
)

func init() {
	if err := gdb.Register(`mysql`, New()); err != nil {
		panic(err)
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for MySQL.
func New() gdb.Driver {
	return &gdb.DriverMysql{}
}
