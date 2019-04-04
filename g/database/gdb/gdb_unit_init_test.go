// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
    "fmt"
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/container/garray"
    "github.com/gogf/gf/g/database/gdb"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/test/gtest"
	"os"
)

const (
    // 初始化表数据量
    INIT_DATA_SIZE = 10
)

var (
	// 数据库对象/接口
	db gdb.DB
)

// 初始化连接参数。
// 测试前需要修改连接参数。
func init() {
	node := gdb.ConfigNode{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Pass:     "",
		Name:     "",
		Type:     "mysql",
		Role:     "master",
		Charset:  "utf8",
		Priority: 1,
	}
	hostname, _ := os.Hostname()
	// 本地测试hack
	if hostname == "ijohn" {
		node.Pass = "12345678"
	}
	gdb.AddConfigNode("test",                 node)
	gdb.AddConfigNode(gdb.DEFAULT_GROUP_NAME, node)
	if r, err := gdb.New(); err != nil {
        gtest.Fatal(err)
	} else {
		db = r
	}
	// 准备测试数据结构
    if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS `test` CHARACTER SET UTF8"); err != nil {
        gtest.Fatal(err)
    }
    // 选择操作数据库
    db.SetSchema("test")
    // 创建默认用户表
    createTable("user")
}

// 创建指定名称的user测试表，当table为空时，创建随机的表名。
// 创建的测试表默认没有任何数据。
// 执行完成后返回该表名。
// TODO 支持更多数据库
func createTable(table...string) (name string) {
    if len(table) > 0 {
        name = table[0]
    } else {
        name = fmt.Sprintf(`user_%d`, gtime.Nanosecond())
    }
    dropTable(name)
    if _, err := db.Exec(fmt.Sprintf(`
    CREATE TABLE %s (
        id int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
        passport varchar(45) NOT NULL COMMENT '账号',
        password char(32) NOT NULL COMMENT '密码',
        nickname varchar(45) NOT NULL COMMENT '昵称',
        create_time timestamp NOT NULL COMMENT '创建时间/注册时间',
        PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, name)); err != nil {
        gtest.Fatal(err)
    }
    return
}

// 删除指定表.
func dropTable(table string) {
    if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
        gtest.Fatal(err)
    }
}

// See createTable.
// 创建测试表，并初始化默认数据。
func createInitTable(table...string) (name string) {
    name   = createTable(table...)
    array := garray.New(true)
    for i := 1; i <= INIT_DATA_SIZE; i++ {
        array.Append(g.Map{
            "id"          : i,
            "passport"    : fmt.Sprintf(`t%d`, i),
            "password"    : fmt.Sprintf(`p%d`, i),
            "nickname"    : fmt.Sprintf(`T%d`, i),
            "create_time" : gtime.Now().String(),
        })
    }
    result, err := db.Table(name).Data(array.Slice()).Insert()
    gtest.Assert(err, nil)

    n, e := result.RowsAffected()
    gtest.Assert(e, nil)
    gtest.Assert(n, INIT_DATA_SIZE)
    return
}
