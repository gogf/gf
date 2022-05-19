// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	_ "github.com/sijms/go-ora/v2"
	"strings"
)

var (
	db     gdb.DB
	dblink gdb.DB
	dbErr  gdb.DB
	ctx    context.Context
)

const (
	TableSize        = 10
	TableName        = "t_user"
	TestSchema1      = "test1"
	TestSchema2      = "test2"
	TableNamePrefix1 = "gf_"
	TestSchema       = "XE"
)

const (
	TestDbIP   = "127.0.0.1"
	TestDbPort = "1521"
	TestDbUser = "system"
	TestDbPass = "oracle"
	TestDbName = "XE"
	TestDbType = "oracle"
)

func init() {
	node := gdb.ConfigNode{
		Host:             TestDbIP,
		Port:             TestDbPort,
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             TestDbName,
		Type:             TestDbType,
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
	}

	nodeLink := gdb.ConfigNode{
		Type: TestDbType,
		Name: TestDbName,
		Link: fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			TestDbType, TestDbUser, TestDbPass, TestDbIP, TestDbPort, TestDbName),
	}

	nodeErr := gdb.ConfigNode{
		Host:    TestDbIP,
		Port:    TestDbPort,
		User:    TestDbUser,
		Pass:    "1234",
		Name:    TestDbName,
		Type:    TestDbType,
		Role:    "master",
		Charset: "utf8",
		Weight:  1,
		Debug:   true,
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

func createTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("user_%d", gtime.Timestamp())
	}

	dropTable(name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		ID NUMBER(10) NOT NULL,
		PASSPORT VARCHAR(45) NOT NULL,
		PASSWORD CHAR(32) NOT NULL,
		NICKNAME VARCHAR(45) NOT NULL,
		CREATE_TIME varchar(45),
		PRIMARY KEY (ID))
	`, name)); err != nil {
		gtest.Fatal(err)
	}

	//db.Schema("test")
	return
}

func createInitTable(table ...string) (name string) {
	name = createTable(table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.Now().String(),
		})
	}
	result, err := db.Insert(context.Background(), name, array.Slice())
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func dropTable(table string) {
	count, err := db.GetCount(ctx, "SELECT COUNT(*) FROM USER_TABLES WHERE TABLE_NAME = ?", strings.ToUpper(table))
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
