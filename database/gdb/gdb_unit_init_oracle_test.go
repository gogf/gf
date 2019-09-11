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
	"strings"
)

var (
	// 数据库对象/接口
	oradb gdb.DB
)

func InitOracle() {
	node := gdb.ConfigNode{
		Host:             "192.168.146.0",
		Port:             "1521",
		User:             "scott",
		Pass:             "tiger",
		Name:             "orcl",
		Type:             "oracle",
		Role:             "master",
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifetime:  600,
	}

	node1 := node
	node1.LinkInfo = fmt.Sprintf("%s/%s@%s", node.User, node.Pass, node.Name)

	gdb.AddConfigNode("oracleNode", node)
	gdb.AddConfigNode("oracleNode", node)
	if r, err := gdb.New("oracleNode"); err != nil {
		gtest.Fatal(err)
	} else {
		oradb = r
	}

	// 创建默认用户表
	createTableOracle("t_user")
}

func createTableOracle(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("t_user_%d", gtime.Nanosecond())
	}

	dropTableOracle(name)

	if _, err := oradb.Exec(fmt.Sprintf(`
	CREATE TABLE %s (
		ID NUMBER(10) NOT NULL,
		PASSPORT VARCHAR(45) NOT NULL,
		PASSWORD CHAR(32) NOT NULL,
		NICKNAME VARCHAR(45) NOT NULL,
		CREATE_TIME varchar(45) NOT NULL,
		PRIMARY KEY (ID))
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return
}

func createInitTableOracle(table ...string) (name string) {
	name = createTableOracle(table...)
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
	result, err := oradb.Table(name).Data(array.Slice()).Insert()
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, INIT_DATA_SIZE)
	return
}

// 删除指定表.
func dropTableOracle(table string) {

	count, err := oradb.GetCount("SELECT COUNT(*) FROM USER_TABLES WHERE TABLE_NAME = ?", strings.ToUpper(table))
	if err != nil {
		gtest.Fatal(err)
	}

	if count == 0 {
		return
	}
	if _, err := oradb.Exec(fmt.Sprintf("DROP TABLE %s", table)); err != nil {
		gtest.Fatal(err)
	}
}
