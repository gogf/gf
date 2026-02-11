// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package mariadb implements gdb.Driver, which supports operations for database MariaDB.
package mariadb

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
)

// Driver is the driver for MariaDB database.
//
// MariaDB is a community-developed, commercially supported fork of the MySQL relational database.
// This driver uses the MySQL protocol to communicate with MariaDB database, as MariaDB maintains
// high compatibility with MySQL protocol.
//
// Although MariaDB is compatible with MySQL protocol, it is packaged as a separate driver component
// rather than reusing the mysql adapter directly. This design allows for future extensibility,
// such as implementing MariaDB-specific features or optimizations.
type Driver struct {
	*mysql.Driver
}

func init() {
	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"mariadb"}
	)
	for _, driverName := range driverNames {
		if err = gdb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
}

// New creates and returns a driver that implements gdb.Driver, which supports operations for MariaDB.
func New() gdb.Driver {
	mysqlDriver := mysql.New().(*mysql.Driver)
	return &Driver{
		Driver: mysqlDriver,
	}
}

// New creates and returns a database object for MariaDB.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	mysqlDB, err := d.Driver.New(core, node)
	if err != nil {
		return nil, err
	}
	return &Driver{
		Driver: mysqlDB.(*mysql.Driver),
	}, nil
}
