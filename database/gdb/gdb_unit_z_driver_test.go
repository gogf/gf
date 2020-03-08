// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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
// gdb.DriverMysql and overwrites its function HandleSqlBeforeExec.
// So if there's any sql execution, it goes through MyDriver.HandleSqlBeforeExec firstly and
// then gdb.DriverMysql.HandleSqlBeforeExec.
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

// HandleSqlBeforeExec handles the sql before posts it to database.
// It here overwrites the same method of gdb.DriverMysql and makes some custom changes.
func (d *MyDriver) HandleSqlBeforeExec(sql string) string {
	latestSqlString.Set(sql)
	return d.DriverMysql.HandleSqlBeforeExec(sql)
}

func init() {
	// It here registers my custom driver in package initialization function "init".
	// You can later using this type in the configuration.
	gdb.Register(customDriverName, &MyDriver{})
}

func Test_Custom_Driver(t *testing.T) {
	gdb.AddConfigNode("driver-test", gdb.ConfigNode{
		Host:    "127.0.0.1",
		Port:    "3306",
		User:    "root",
		Pass:    "12345678",
		Name:    "test",
		Type:    customDriverName,
		Role:    "master",
		Charset: "utf8",
	})
	gtest.Assert(latestSqlString.Val(), "")
	sqlString := "select 10000"
	value, err := g.DB("driver-test").GetValue(sqlString)
	gtest.Assert(err, nil)
	gtest.Assert(value, 10000)
	gtest.Assert(latestSqlString.Val(), sqlString)
}
