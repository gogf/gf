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
	pgdb gdb.DB
)

// 初始化连接参数。
// 测试前需要修改连接参数。
func InitPgsql() {
	node := gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "5432",
		User:             "postgres",
		Pass:             "password",
		Name:             "travis_ci_test",
		Type:             "pgsql",
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifetime:  600,
	}

	gdb.AddConfigNode("pgsqlNode", node)
	gdb.AddConfigNode("pgsqlNode", node)
	if r, err := gdb.New("pgsqlNode"); err != nil {
		gtest.Fatal(err)
	} else {
		pgdb = r
	}

	/*if _, err := pgdb.Exec(fmt.Sprintf("drop database if exists %s", SCHEMA1)); err != nil {
		gtest.Error(err)
	}
	schemaTemplate := "CREATE DATABASE %s"
	if _, err := pgdb.Exec(fmt.Sprintf(schemaTemplate, SCHEMA1)); err != nil {
		gtest.Error(err)
	}
	if _, err := pgdb.Exec(fmt.Sprintf("SET search_path TO %s", SCHEMA1)); err != nil {
		gtest.Error(err)
	}
	pgdb.SetSchema(SCHEMA1)
	*/

	// 创建默认用户表
	createTablePgsql("t_user")

}

// 创建指定名称的user测试表，当table为空时，创建随机的表名。
// 创建的测试表默认没有任何数据。
// 执行完成后返回该表名。
// TODO 支持更多数据库
func createTablePgsql(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`user_%d`, gtime.Nanosecond())
	}
	dropTablePgsql(name)
	if _, err := pgdb.Exec(fmt.Sprintf(`
    CREATE TABLE %s (
        id bigint  NOT NULL,
        passport varchar(45),
        password char(32) NOT NULL,
        nickname varchar(45) NOT NULL,
        create_time timestamp NOT NULL,
        PRIMARY KEY (id)
    ) ;
    `, name)); err != nil {
		gtest.Fatal(err)
	}

	return
}

// 创建测试表，并初始化默认数据。
func createInitTablePgsql(table ...string) (name string) {
	name = createTablePgsql(table...)
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
	result, err := pgdb.Table(name).Data(array.Slice()).Insert()
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, INIT_DATA_SIZE)
	return
}

// 删除指定表.
func dropTablePgsql(table string) {
	if _, err := pgdb.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
		gtest.Fatal(err)
	}
}
