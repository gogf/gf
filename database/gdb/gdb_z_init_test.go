// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	TableSize        = 10
	TableName        = "user"
	TestSchema1      = "test1"
	TestSchema2      = "test2"
	TableNamePrefix1 = "gf_"
	TestDbUser       = "fanwei"
	TestDbPass       = "fw123456"
	CreateTime       = "2018-10-24 10:00:00"
)

var (
	db         gdb.DB
	dbPrefix   gdb.DB
	dbInvalid  gdb.DB
	configNode gdb.ConfigNode
	ctx        = context.TODO()
)

func init() {
	parser, err := gcmd.Parse(g.MapStrBool{
		"name": true,
		"type": true,
	}, false)
	gtest.AssertNil(err)
	configNode = gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "3306",
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             parser.GetOpt("name", "").String(),
		Type:             parser.GetOpt("type", "mysql").String(),
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifeTime:  600,
	}
	nodePrefix := configNode
	nodePrefix.Prefix = TableNamePrefix1

	nodeInvalid := configNode
	nodeInvalid.Port = "3307"

	gdb.AddConfigNode("test", configNode)
	gdb.AddConfigNode("prefix", nodePrefix)
	gdb.AddConfigNode("nodeinvalid", nodeInvalid)
	gdb.AddConfigNode(gdb.DefaultGroupName, configNode)

	// Default db.
	if r, err := gdb.New(); err != nil {
		gtest.Error(err)
	} else {
		db = r
	}
	schemaTemplate := "CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET UTF8"
	if _, err := db.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema1)); err != nil {
		gtest.Error(err)
	}
	if _, err := db.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema2)); err != nil {
		gtest.Error(err)
	}
	db.SetSchema(TestSchema1)

	// Prefix db.
	if r, err := gdb.New("prefix"); err != nil {
		gtest.Error(err)
	} else {
		dbPrefix = r
	}
	if _, err := dbPrefix.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema1)); err != nil {
		gtest.Error(err)
	}
	if _, err := dbPrefix.Exec(ctx, fmt.Sprintf(schemaTemplate, TestSchema2)); err != nil {
		gtest.Error(err)
	}
	dbPrefix.SetSchema(TestSchema1)

	// Invalid db.
	if r, err := gdb.New("nodeinvalid"); err != nil {
		gtest.Error(err)
	} else {
		dbInvalid = r
	}
	dbInvalid.SetSchema(TestSchema1)
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

func createTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	}
	dropTableWithDb(db, name)

	switch configNode.Type {
	case "sqlite":
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   id bigint unsigned NOT NULL AUTO_INCREMENT,
		   passport varchar(45),
		   password char(32) NOT NULL,
		   nickname varchar(45) NOT NULL,
		   create_time timestamp NOT NULL,
		   PRIMARY KEY (id)
		) ;`, name,
		)); err != nil {
			gtest.Fatal(err)
		}
	case "pgsql":
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   id bigint  NOT NULL,
		   passport varchar(45),
		   password char(32) NOT NULL,
		   nickname varchar(45) NOT NULL,
		   create_time timestamp NOT NULL,
		   PRIMARY KEY (id)
		) ;`, name,
		)); err != nil {
			gtest.Fatal(err)
		}
	case "mssql":
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		CREATE TABLE %s (
		ID numeric(10,0) NOT NULL,
		PASSPORT VARCHAR(45) NOT NULL,
		PASSWORD CHAR(32) NOT NULL,
		NICKNAME VARCHAR(45) NOT NULL,
		CREATE_TIME datetime NOT NULL,
		PRIMARY KEY (ID))`,
			name, name,
		)); err != nil {
			gtest.Fatal(err)
		}
	case "oracle":
		if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		ID NUMBER(10) NOT NULL,
		PASSPORT VARCHAR(45) NOT NULL,
		PASSWORD CHAR(32) NOT NULL,
		NICKNAME VARCHAR(45) NOT NULL,
		CREATE_TIME varchar(45) NOT NULL,
		PRIMARY KEY (ID))
		`, name,
		)); err != nil {
			gtest.Fatal(err)
		}
	case "mysql":
		if _, err := db.Exec(ctx, fmt.Sprintf(`
	    CREATE TABLE %s (
	        id          int(10) unsigned NOT NULL AUTO_INCREMENT,
	        passport    varchar(45) NULL,
	        password    char(32) NULL,
	        nickname    varchar(45) NULL,
	        create_time timestamp NULL,
	        PRIMARY KEY (id)
	    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	    `, name,
		)); err != nil {
			gtest.Fatal(err)
		}
	}

	return
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
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
}
