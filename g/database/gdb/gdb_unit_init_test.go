// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"os"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/container/garray"

	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
)

const (
	INIT_DATA_SIZE = 10      // 初始化表数据量
	TABLE          = "user"  // 测试数据表
	SCHEMA1        = "test1" // 测试数据库1
	SCHEMA2        = "test2" // 测试数据库2
)

var (
	// 测试包变量，ORM对象
	db gdb.DB
)

// 初始化连接参数。
// 测试前需要修改连接参数。
func init() {
	node := gdb.ConfigNode{
		Host:    "127.0.0.1",
		Port:    "3306",
		User:    "root",
		Pass:    "",
		Name:    "",
		Type:    "mysql",
		Role:    "master",
		Charset: "utf8",
		Weight:  1,
	}
	// 作者本地测试hack
	if hostname, _ := os.Hostname(); hostname == "ijohn" {
		node.Pass = "12345678"
	}
	gdb.AddConfigNode("test", node)
	gdb.AddConfigNode(gdb.DEFAULT_GROUP_NAME, node)
	if r, err := gdb.New(); err != nil {
		gtest.Error(err)
	} else {
		db = r
	}
	// 准备测试数据结构：数据库
	schemaTemplate := "CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET UTF8"
	if _, err := db.Exec(fmt.Sprintf(schemaTemplate, SCHEMA1)); err != nil {
		gtest.Error(err)
	}
	// 多个数据库，用于测试数据库切换
	if _, err := db.Exec(fmt.Sprintf(schemaTemplate, SCHEMA2)); err != nil {
		gtest.Error(err)
	}
	// 设置默认操作数据库
	db.SetSchema(SCHEMA1)
	// 创建默认用户表
	createTable(TABLE)
}

// 创建指定名称的user测试表，当table为空时，创建随机的表名。
// 创建的测试表默认没有任何数据。
// 执行完成后返回该表名。
func createTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TABLE, gtime.Nanosecond())
	}
	dropTable(name)
	if _, err := db.Exec(fmt.Sprintf(`
    CREATE TABLE %s (
        id          int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
        passport    varchar(45) NOT NULL COMMENT '账号',
        password    char(32) NOT NULL COMMENT '密码',
        nickname    varchar(45) NOT NULL COMMENT '昵称',
        create_time timestamp NOT NULL COMMENT '创建时间/注册时间',
        PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, name)); err != nil {
		gtest.Error(err)
	}
	return
}

// 创建测试表，并初始化默认数据。
func createInitTable(table ...string) (name string) {
	name = createTable(table...)
	array := garray.New(true)
	for i := 1; i <= INIT_DATA_SIZE; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),
		})
	}
	result, err := db.Table(name).Data(array.Slice()).Insert()
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, INIT_DATA_SIZE)
	return
}

// 删除指定表.
func dropTable(table string) {
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
}
