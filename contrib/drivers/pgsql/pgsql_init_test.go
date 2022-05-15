// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"context"
	"fmt"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/lib/pq"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	TableSize   = 10
	TablePrefix = "t_"
	TableName   = "test"
	TestSchema1 = "test1"
	TestSchema2 = "test2"
	TestDbUser  = "root"
	TestDbPass  = "12345678"
	CreateTime  = "2018-10-24 10:00:00"
)

var (
	db         gdb.DB
	dbT1       gdb.DB
	dbT2       gdb.DB
	configNode gdb.ConfigNode
	ctx        = context.TODO()
)

func init() {
	configNode = gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "5432",
		User:             TestDbUser,
		Pass:             TestDbPass,
		Timezone:         "Asia/Shanghai", // For calculating UT cases of datetime zones in convenience.
		Type:             "pgsql",
		Name:             TableName,
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifeTime:  600,
	}

	subNode1 := configNode
	subNode1.Role = "slave"
	subNode1.Name = TestSchema1

	subNode2 := configNode
	subNode2.Role = "slave"
	subNode2.Name = TestSchema2

	//pgsql only permit to connect to the designation database.
	//so you need to create the pgsql database before you use orm

	gdb.AddConfigNode(gdb.DefaultGroupName, configNode)
	if r, err := gdb.New(configNode); err != nil {
		gtest.Fatal(err)
	} else {
		db = r
	}

	nodeT1 := gdb.ConfigNode{
		Type: "pgsql",
		Link: fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			configNode.User, configNode.Pass, configNode.Host, configNode.Port, TestSchema1),
	}

	nodeT2 := gdb.ConfigNode{
		Type: "pgsql",
		Link: fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			configNode.User, "configNode.Pass", configNode.Host, configNode.Port, TestSchema2),
	}

	gdb.AddConfigNode("dbTest", nodeT1)
	if r, err := gdb.New(nodeT1); err != nil {
		gtest.Fatal(err)
	} else {
		dbT1 = r
	}

	gdb.AddConfigNode("dbTest", nodeT2)
	if r, err := gdb.New(nodeT2); err != nil {
		gtest.Fatal(err)
	} else {
		dbT2 = r
	}
}

func createTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TablePrefix+"test", gtime.TimestampNano())
	}

	dropTableWithDb(db, name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   id bigint  NOT NULL,
		   passport varchar(45),
		   password varchar(32) NOT NULL,
		   nickname varchar(45) NOT NULL,
		   create_time timestamp NOT NULL,
		   PRIMARY KEY (id)
		) ;`, name,
	)); err != nil {
		gtest.Fatal(err)
	}
	return
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
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
		gtest.Error(err)
	}
}
