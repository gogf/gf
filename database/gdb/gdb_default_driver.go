// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
)

// DriverTest is the driver for mysql database.
type DriverTest struct {
	*Core
}

func init() {
	if err := Register("test", &DriverTest{}); err != nil {
		panic(err)
	}
}

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverTest) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverTest{
		Core: core,
	}, nil
}

// Open creates and returns an underlying sql.DB object for mysql.
// Note that it converts time.Time argument to local timezone in default.
func (d *DriverTest) Open(config *ConfigNode) (db *sql.DB, err error) {
	return
}

// PingMaster pings the master node to check authentication or keeps the connection alive.
func (d *DriverTest) PingMaster() error {
	return nil
}

// PingSlave pings the slave node to check authentication or keeps the connection alive.
func (d *DriverTest) PingSlave() error {
	return nil
}
