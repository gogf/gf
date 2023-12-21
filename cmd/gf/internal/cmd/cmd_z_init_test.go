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
	ctx    = context.Background()
	testDB gdb.DB
	link   = "mysql:root:12345678@tcp(127.0.0.1:3306)/test?loc=Local&parseTime=true"
)

func init() {
	var err error
	testDB, err = gdb.New(gdb.ConfigNode{
		Link: link,
	})
	if err != nil {
		panic(err)
	}
}

func dropTableWithDb(db gdb.DB, table string) {
	dropTableStmt := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)
	if _, err := db.Exec(ctx, dropTableStmt); err != nil {
		gtest.Error(err)
	}
}
