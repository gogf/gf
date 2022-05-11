// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	// 数据库对象/接口
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
	TestDbUser       = "sa"
	TestDbPass       = "LoremIpsum86"
	CreateTime       = "2018-10-24 10:00:00"
)

func init() {
	node := gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "1433",
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             "test",
		Type:             "mssql",
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
	}

	nodeLink := gdb.ConfigNode{
		Type: "mssql",
		Name: "test",
		Link: fmt.Sprintf("user id=%s;password=%s;server=%s;port=%s;database=%s;encrypt=disable",
			node.User, node.Pass, node.Host, node.Port, node.Name),
	}

	nodeErr := gdb.ConfigNode{
		Type: "mssql",
		Link: fmt.Sprintf("user id=%s;password=%s;server=%s;port=%s;database=%s;encrypt=disable",
			node.User, "node.Pass", node.Host, node.Port, node.Name),
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

	if _, err := db.Exec(context.Background(), fmt.Sprintf(`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		CREATE TABLE %s (
		ID numeric(10,0) NOT NULL,
		PASSPORT VARCHAR(45)  NULL,
		PASSWORD VARCHAR(32)  NULL,
		NICKNAME VARCHAR(45)  NULL,
		CREATE_TIME datetime NULL,
		PRIMARY KEY (ID))
	`, name, name)); err != nil {
		gtest.Fatal(err)
	}

	// 选择操作数据库
	db.Schema("test")

	//db.SetDebug(true)
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

// 删除指定表.
func dropTable(table string) {
	if _, err := db.Exec(context.Background(), fmt.Sprintf(`
		IF EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		DROP TABLE %s
	`, table, table)); err != nil {
		gtest.Fatal(err)
	}
}
