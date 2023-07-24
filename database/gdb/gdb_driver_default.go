// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
)

// DriverDefault is the default driver for mysql database, which does nothing.
type DriverDefault struct {
	*Core
}

func init() {
	if err := Register("default", &DriverDefault{}); err != nil {
		panic(err)
	}
}

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverDefault) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverDefault{
		Core: core,
	}, nil
}

// Open creates and returns an underlying sql.DB object for mysql.
// Note that it converts time.Time argument to local timezone in default.
func (d *DriverDefault) Open(config *ConfigNode) (db *sql.DB, err error) {
	return
}

// PingMaster pings the master node to check authentication or keeps the connection alive.
func (d *DriverDefault) PingMaster() error {
	return nil
}

// PingSlave pings the slave node to check authentication or keeps the connection alive.
func (d *DriverDefault) PingSlave() error {
	return nil
}
