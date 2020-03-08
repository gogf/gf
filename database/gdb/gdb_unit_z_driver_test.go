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
type MyDriver struct {
	*gdb.DriverMysql
}

var (
	customDriverName = "MyDriver"
	lastSqlString    = gtype.NewString() // For unit testing only.
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
	lastSqlString.Set(sql)
	return d.DriverMysql.HandleSqlBeforeExec(sql)
}

func init() {
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
	gtest.Assert(lastSqlString.Val(), "")
	sqlString := "select 10000"
	value, err := g.DB("driver-test").GetValue(sqlString)
	gtest.Assert(err, nil)
	gtest.Assert(value, 10000)
	gtest.Assert(lastSqlString.Val(), sqlString)
}
