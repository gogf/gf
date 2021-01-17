// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

const (
	SIZE    = 10
	TABLE   = "user"
	SCHEMA1 = "test1"
	SCHEMA2 = "test2"
	PREFIX1 = "gf_"
	USER    = "root"
	PASS    = "12345678"
)

var (
	db         gdb.DB
	dbPrefix   gdb.DB
	configNode gdb.ConfigNode
)

func init() {
	parser, err := gcmd.Parse(g.MapStrBool{
		"name": true,
		"type": true,
	}, false)
	gtest.Assert(err, nil)
	configNode = gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "3306",
		User:             USER,
		Pass:             PASS,
		Name:             parser.GetOpt("name", ""),
		Type:             parser.GetOpt("type", "mysql"),
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifetime:  600,
	}
	nodePrefix := configNode
	nodePrefix.Prefix = PREFIX1
	gdb.AddConfigNode("test", configNode)
	gdb.AddConfigNode("prefix", nodePrefix)
	gdb.AddConfigNode(gdb.DefaultGroupName, configNode)
	// Default db.
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

	// Prefix db.
	if r, err := gdb.New("prefix"); err != nil {
		gtest.Error(err)
	} else {
		dbPrefix = r
	}
	if _, err := dbPrefix.Exec(fmt.Sprintf(schemaTemplate, SCHEMA1)); err != nil {
		gtest.Error(err)
	}
	if _, err := dbPrefix.Exec(fmt.Sprintf(schemaTemplate, SCHEMA2)); err != nil {
		gtest.Error(err)
	}
	dbPrefix.SetSchema(SCHEMA1)
	createTable(TABLE)
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
		name = fmt.Sprintf(`%s_%d`, TABLE, gtime.TimestampNano())
	}
	dropTableWithDb(db, name)

	switch configNode.Type {
	case "sqlite":
		if _, err := db.Exec(fmt.Sprintf(`
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
		if _, err := db.Exec(fmt.Sprintf(`
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
		if _, err := db.Exec(fmt.Sprintf(`
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
		if _, err := db.Exec(fmt.Sprintf(`
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
		if _, err := db.Exec(fmt.Sprintf(`
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
	for i := 1; i <= SIZE; i++ {
		array.Append(g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr("2018-10-24 10:00:00").String(),
		})
	}

	result, err := db.BatchInsert(name, array.Slice())
	gtest.Assert(err, nil)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, SIZE)
	return
}

func dropTableWithDb(db gdb.DB, table string) {
	if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
}
