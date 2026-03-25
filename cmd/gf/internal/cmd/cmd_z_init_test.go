// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	ctx      = context.Background()
	testDB   gdb.DB
	testPgDB gdb.DB
	link     = "mysql:root:12345678@tcp(127.0.0.1:3306)/test?loc=Local&parseTime=true"
	linkPg   = "pgsql:postgres:12345678@tcp(127.0.0.1:5432)/test"
)

func init() {
	var err error
	testDB, err = gdb.New(gdb.ConfigNode{
		Link: link,
	})
	if err != nil {
		panic(err)
	}
	// PostgreSQL connection (optional, may not be available in all environments)
	testPgDB, _ = gdb.New(gdb.ConfigNode{
		Link: linkPg,
	})
}

func dropTableWithDb(db gdb.DB, table string) {
	dropTableStmt := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
	if _, err := db.Exec(ctx, dropTableStmt); err != nil {
		gtest.Error(err)
	}
}

// dropTableStd uses standard SQL syntax compatible with MySQL and PostgreSQL.
func dropTableStd(db gdb.DB, table string) {
	dropTableStmt := fmt.Sprintf("DROP TABLE IF EXISTS %s", table)
	if _, err := db.Exec(ctx, dropTableStmt); err != nil {
		gtest.Error(err)
	}
}
