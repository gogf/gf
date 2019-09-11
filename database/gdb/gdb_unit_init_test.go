// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	// _ "github.com/denisenkom/go-mssqldb"
	// _ "github.com/lib/pq"
	// _ "github.com/mattn/go-oci8"
	_ "github.com/gogf/gf/database/gdb"
)

func init() {
	//InitPgsql()
	//InitOracle()
	//InitMssql()
	InitMysql()
}
