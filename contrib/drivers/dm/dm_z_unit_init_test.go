// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	db        gdb.DB
	dblink    gdb.DB
	dbErr     gdb.DB
	ctx       context.Context
	TableSize = 10
)

const (
	TestDBHost  = "127.0.0.1"
	TestDBPort  = "5236"
	TestDBUser  = "SYSDBA"
	TestDBPass  = "SYSDBA001"
	TestDBName  = "SYSDBA"
	TestDBType  = "dm"
	TestCharset = "utf8"
)

type User struct {
	ID          int64     `orm:"id"`
	AccountName string    `orm:"account_name"`
	PwdReset    int64     `orm:"pwd_reset"`
	AttrIndex   int64     `orm:"attr_index"`
	Enabled     int64     `orm:"enabled"`
	Deleted     int64     `orm:"deleted"`
	CreatedBy   string    `orm:"created_by"`
	CreatedTime time.Time `orm:"created_time"`
	UpdatedBy   string    `orm:"updated_by"`
	UpdatedTime time.Time `orm:"updated_time"`
}

func init() {
	node := gdb.ConfigNode{
		Host:             TestDBHost,
		Port:             TestDBPort,
		User:             TestDBUser,
		Pass:             TestDBPass,
		Name:             TestDBName,
		Type:             TestDBType,
		Role:             "master",
		Charset:          TestCharset,
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		CreatedAt:        "created_time",
		UpdatedAt:        "updated_time",
	}

	// todo
	nodeLink := gdb.ConfigNode{
		Type: TestDBType,
		Name: TestDBName,
		Link: fmt.Sprintf(
			"dm:%s:%s@tcp(%s:%s)/%s?charset=%s",
			TestDBUser, TestDBPass, TestDBHost, TestDBPort, TestDBName, TestCharset,
		),
	}

	nodeErr := gdb.ConfigNode{
		Host:    TestDBHost,
		Port:    TestDBPort,
		User:    TestDBUser,
		Pass:    "1234",
		Name:    TestDBName,
		Type:    TestDBType,
		Role:    "master",
		Charset: TestCharset,
		Weight:  1,
	}

	gdb.AddConfigNode(gdb.DefaultGroupName, node)
	if r, err := gdb.New(node); err != nil {
		gtest.Fatal(err)
	} else {
		db = r
	}

	gdb.AddConfigNode("dblink", nodeLink)
	if r, err := gdb.New(nodeLink); err != nil {
		gtest.Fatal(err)
	} else {
		dblink = r
	}

	gdb.AddConfigNode("dbErr", nodeErr)
	if r, err := gdb.New(nodeErr); err != nil {
		gtest.Fatal(err)
	} else {
		dbErr = r
	}

	ctx = context.Background()
}

func dropTable(table string) {
	count, err := db.GetCount(
		ctx,
		"SELECT COUNT(*) FROM all_tables WHERE owner = ? And table_name= ?", TestDBName, strings.ToUpper(table),
	)
	if err != nil {
		gtest.Fatal(err)
	}

	if count == 0 {
		return
	}
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE %s", table)); err != nil {
		gtest.Fatal(err)
	}
}

func createTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("random_%d", gtime.Timestamp())
	}

	dropTable(name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE "%s"
(
"ID" BIGINT NOT NULL,
"ACCOUNT_NAME" VARCHAR(128) DEFAULT '' NOT NULL,
"PWD_RESET" TINYINT DEFAULT 0 NOT NULL,
"ENABLED" INT DEFAULT 1 NOT NULL,
"DELETED" INT DEFAULT 0 NOT NULL,
"ATTR_INDEX" INT DEFAULT 0 ,
"CREATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"CREATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
"UPDATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"UPDATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
NOT CLUSTER PRIMARY KEY("ID")) STORAGE(ON "MAIN", CLUSTERBTR) ;
	`, name)); err != nil {
		gtest.Fatal(err)
	}

	return
}

func createInitTable(table ...string) (name string) {
	name = createTable(table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":           i,
			"account_name": fmt.Sprintf(`name_%d`, i),
			"pwd_reset":    0,
			"attr_index":   i,
			"create_time":  gtime.Now().String(),
		})
	}
	result, err := db.Schema(TestDBName).Insert(context.Background(), name, array.Slice())
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func createTableFalse(table ...string) (name string, err error) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("random_%d", gtime.Timestamp())
	}

	dropTable(name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE "%s"
(
"ID" BIGINT NOT NULL,
"ACCOUNT_NAME" VARCHAR(128) DEFAULT '' NOT NULL,
"PWD_RESET" TINYINT DEFAULT 0 NOT NULL,
"ENABLED" INT DEFAULT 1 NOT NULL,
"DELETED" INT DEFAULT 0 NOT NULL,
"INDEX" INT DEFAULT 0 ,
"ATTR_INDEX" INT DEFAULT 0 ,
"CREATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"CREATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
"UPDATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"UPDATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
NOT CLUSTER PRIMARY KEY("ID")) STORAGE(ON "MAIN", CLUSTERBTR) ;
	`, name)); err != nil {
		// gtest.Fatal(err)
		return name, fmt.Errorf("createTableFalse")
	}

	return name, nil
}
