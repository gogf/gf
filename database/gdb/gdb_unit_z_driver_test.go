// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

// MyDriver is a custom database driver, which is used for testing only.
type MyDriver struct {
	*gdb.Core
}

var (
	myCustomDriverName = "mydriver"
	lastSqlString      = gtype.NewString() // For unit testing only.
)

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *MyDriver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &MyDriver{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for mysql.
func (d *MyDriver) Open(config *gdb.ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true&parseTime=true&loc=Local",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
	}
	// It uses mysql driver as underlying sql driver.
	if db, err := sql.Open("mysql", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// GetChars returns the security char for this type of database.
func (d *MyDriver) GetChars() (charLeft string, charRight string) {
	return "`", "`"
}

// HandleSqlBeforeExec handles the sql before posts it to database.
func (d *MyDriver) HandleSqlBeforeExec(sql string) string {
	lastSqlString.Set(sql)
	return sql
}

func init() {
	gdb.Register(myCustomDriverName, &MyDriver{})
}

func Test_Custom_Driver(t *testing.T) {
	gdb.AddConfigNode("driver-test", gdb.ConfigNode{
		Host:    "127.0.0.1",
		Port:    "3306",
		User:    "root",
		Pass:    "12345678",
		Name:    "test",
		Type:    myCustomDriverName,
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
