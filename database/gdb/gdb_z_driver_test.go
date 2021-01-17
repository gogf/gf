// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

// MyDriver is a custom database driver, which is used for testing only.
// For simplifying the unit testing case purpose, MyDriver struct inherits the mysql driver
// gdb.DriverMysql and overwrites its function HandleSqlBeforeCommit.
// So if there's any sql execution, it goes through MyDriver.HandleSqlBeforeCommit firstly and
// then gdb.DriverMysql.HandleSqlBeforeCommit.
// You can call it sql "HOOK" or "HiJack" as your will.
type MyDriver struct {
	*gdb.DriverMysql
}

var (
	customDriverName = "MyDriver"
	latestSqlString  = gtype.NewString() // For simplifying unit testing only.
)

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *MyDriver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &MyDriver{
		&gdb.DriverMysql{
			Core: core,
		},
	}, nil
}

// HandleSqlBeforeCommit handles the sql before posts it to database.
// It here overwrites the same method of gdb.DriverMysql and makes some custom changes.
func (d *MyDriver) HandleSqlBeforeCommit(link gdb.Link, sql string, args []interface{}) (string, []interface{}) {
	latestSqlString.Set(sql)
	return d.DriverMysql.HandleSqlBeforeCommit(link, sql, args)
}

func init() {
	// It here registers my custom driver in package initialization function "init".
	// You can later use this type in the database configuration.
	gdb.Register(customDriverName, &MyDriver{})
}

func Test_Custom_Driver(t *testing.T) {
	gdb.AddConfigNode("driver-test", gdb.ConfigNode{
		Host:    "127.0.0.1",
		Port:    "3306",
		User:    USER,
		Pass:    PASS,
		Name:    "test",
		Type:    customDriverName,
		Role:    "master",
		Charset: "utf8",
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(latestSqlString.Val(), "")
		sqlString := "select 10000"
		value, err := g.DB("driver-test").GetValue(sqlString)
		t.Assert(err, nil)
		t.Assert(value, 10000)
		t.Assert(latestSqlString.Val(), sqlString)
	})
}
