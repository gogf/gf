// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"

	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	TableSize   = 10
	TablePrefix = "t_"
	SchemaName  = "test"
	CreateTime  = "2018-10-24 10:00:00"
)

var (
	db         gdb.DB
	configNode gdb.ConfigNode
	ctx        = context.TODO()
)

func init() {
	configNode = gdb.ConfigNode{
		Link: `pgsql:postgres:12345678@tcp(127.0.0.1:5432)`,
	}

	//pgsql only permit to connect to the designation database.
	//so you need to create the pgsql database before you use orm
	gdb.AddConfigNode(gdb.DefaultGroupName, configNode)
	if r, err := gdb.New(configNode); err != nil {
		gtest.Fatal(err)
	} else {
		db = r
	}

	if configNode.Name == "" {
		schemaTemplate := "SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')"
		if _, err := db.Exec(ctx, fmt.Sprintf(schemaTemplate, SchemaName, SchemaName)); err != nil {
			gtest.Error(err)
		}

		db = db.Schema(SchemaName)
	} else {
		db = db.Schema(configNode.Name)
	}

}

func createTable(table ...string) string {
	return createTableWithDb(db, table...)
}

func createInitTable(table ...string) string {
	return createInitTableWithDb(db, table...)
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
		   	id bigserial  NOT NULL,
		   	passport varchar(45) NOT NULL,
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
