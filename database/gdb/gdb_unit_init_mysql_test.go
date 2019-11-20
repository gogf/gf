// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

const (
	INIT_DATA_SIZE = 10
	TABLE          = "user"
	SCHEMA1        = "test1"
	SCHEMA2        = "test2"
)

var (
	db gdb.DB
)

func InitMysql() {
	node := gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "3306",
		User:             "root",
		Pass:             "12345678",
		Name:             "",
		Type:             "mysql",
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifetime:  600,
	}
	gdb.AddConfigNode("test", node)
	gdb.AddConfigNode(gdb.DEFAULT_GROUP_NAME, node)
	if r, err := gdb.New(); err != nil {
		gtest.Error(err)
	} else {
		db = r
	}

	schemaTemplate := "CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET UTF8"
	if _, err := db.Exec(fmt.Sprintf(schemaTemplate, SCHEMA1)); err != nil {
		gtest.Error(err)
	}

	if _, err := db.Exec(fmt.Sprintf(schemaTemplate, SCHEMA2)); err != nil {
		gtest.Error(err)
	}

	db.SetSchema(SCHEMA1)

	createTable(TABLE)
}

func createTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TABLE, gtime.Nanosecond())
	}
	dropTable(name)
	if _, err := db.Exec(fmt.Sprintf(`
    CREATE TABLE %s (
        id          int(10) unsigned NOT NULL AUTO_INCREMENT,
        passport    varchar(45) NULL,
        password    char(32) NULL,
        nickname    varchar(45) NULL,
        create_time timestamp NULL,
        PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, name)); err != nil {
		gtest.Error(err)
	}
	return
}

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

func dropTable(table string) {
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
}
