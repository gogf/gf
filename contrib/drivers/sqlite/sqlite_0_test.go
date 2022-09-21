// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite_test

import (
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	db         gdb.DB
	dbPrefix   gdb.DB
	dbInvalid  gdb.DB
	configNode gdb.ConfigNode
	dbDir      = gfile.Temp("sqlite")
	ctx        = gctx.New()

	// Error
	ErrorSave = gerror.NewCode(gcode.CodeNotSupported, `Save operation is not supported by sqlite driver`)
)

const (
	TableSize       = 10
	TableName       = "user"
	TestSchema1     = "test1"
	TestSchema2     = "test2"
	TableNamePrefix = "gf_"
	CreateTime      = "2018-10-24 10:00:00"
	DBGroupTest     = "test"
	DBGroupPrefix   = "prefix"
	DBGroupInvalid  = "invalid"
)

func init() {
	fmt.Println("init sqlite db start")

	if err := gfile.Mkdir(dbDir); err != nil {
		gtest.Error(err)
	}

	fmt.Println("init sqlite db dir: ", dbDir)

	configNode = gdb.ConfigNode{
		Type:    "sqlite",
		Link:    gfile.Join(dbDir, "test.db"),
		Charset: "utf8",
	}
	nodePrefix := configNode
	nodePrefix.Prefix = TableNamePrefix

	nodeInvalid := configNode

	gdb.AddConfigNode(DBGroupTest, configNode)
	gdb.AddConfigNode(DBGroupPrefix, nodePrefix)
	gdb.AddConfigNode(DBGroupInvalid, nodeInvalid)
	gdb.AddConfigNode(gdb.DefaultGroupName, configNode)

	// Default db.
	if r, err := gdb.NewByGroup(); err != nil {
		gtest.Error(err)
	} else {
		db = r
	}

	// Prefix db.
	if r, err := gdb.NewByGroup(DBGroupPrefix); err != nil {
		gtest.Error(err)
	} else {
		dbPrefix = r
	}

	// Invalid db.
	if r, err := gdb.NewByGroup(DBGroupInvalid); err != nil {
		gtest.Error(err)
	} else {
		dbInvalid = r
	}

	fmt.Println("init sqlite db finish")
}

func createTable(table ...string) string {
	return createTableWithDb(db, table...)
}

func createInitTable(table ...string) string {
	return createInitTableWithDb(db, table...)
}

func dropTable(table string) {
	dropTableWithDb(db, table)
}

func createTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	}
	dropTableWithDb(db, name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		id          INTEGER       PRIMARY KEY AUTOINCREMENT
									UNIQUE
									NOT NULL,
		passport    VARCHAR(45)  NOT NULL
									DEFAULT passport,
		password    VARCHAR(128) NOT NULL
									DEFAULT password,
		nickname    VARCHAR(45),
		create_time DATETIME
	);
	`, name,
	)); err != nil {
		gtest.Fatal(err)
	}

	return
}

func createInitTableWithDb(db gdb.DB, table ...string) (name string) {
	name = createTableWithDb(db, table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr(CreateTime).String(),
		})
	}

	result, err := db.Insert(ctx, name, array.Slice())
	gtest.AssertNil(err)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func dropTableWithDb(db gdb.DB, table string) {
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
}
