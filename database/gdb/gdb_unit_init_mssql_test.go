// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

var (
	// 数据库对象/接口
	msdb gdb.DB
)

func InitMssql() {
	node := gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "1433",
		User:             "sa",
		Pass:             "123456",
		Name:             "test",
		Type:             "mssql",
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifetime:  600,
	}

	gdb.AddConfigNode("mssqlGroup", node)
	gdb.AddConfigNode("mssqlGroup", node)
	if r, err := gdb.New("mssqlGroup"); err != nil {
		gtest.Fatal(err)
	} else {
		msdb = r
	}

	// 创建默认用户表
	createTableMssql("t_user")
	//msdb.SetDebug(true)
}

func createTableMssql(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("user_%d", gtime.Nanosecond())
	}

	dropTableMssql(name)

	if _, err := msdb.Exec(fmt.Sprintf(`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		CREATE TABLE %s (
		ID numeric(10,0) NOT NULL,
		PASSPORT VARCHAR(45) NOT NULL,
		PASSWORD CHAR(32) NOT NULL,
		NICKNAME VARCHAR(45) NOT NULL,
		CREATE_TIME datetime NOT NULL,
		PRIMARY KEY (ID))
	`, name, name)); err != nil {
		gtest.Fatal(err)
	}

	//msdb.Exec("DROP DATABASE test")
	//msdb.Exec("CREATE DATABASE test")

	// 选择操作数据库
	msdb.SetSchema("test")

	//msdb.SetDebug(true)
	return
}

func createInitTableMssql(table ...string) (name string) {
	name = createTableMssql(table...)
	array := garray.New(true)
	for i := 1; i <= INIT_DATA_SIZE; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`t%d`, i),
			"password":    fmt.Sprintf(`p%d`, i),
			"nickname":    fmt.Sprintf(`T%d`, i),
			"create_time": gtime.Now().String(),
		})
	}
	result, err := msdb.Table(name).Data(array.Slice()).Insert()
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, INIT_DATA_SIZE)
	return
}

// 删除指定表.
func dropTableMssql(table string) {
	if _, err := msdb.Exec(fmt.Sprintf(`
		IF EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		DROP TABLE %s
	`, table, table)); err != nil {
		gtest.Fatal(err)
	}
}
