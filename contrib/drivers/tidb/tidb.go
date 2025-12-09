// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package tidb implements gdb.Driver, which supports operations for database TiDB.
package tidb

import (
	"github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Driver is the driver for TiDB database.
//
// TiDB is an open-source NewSQL database that supports Hybrid Transactional and Analytical Processing (HTAP).
// This driver uses the MySQL protocol to communicate with TiDB database, as TiDB is designed to be highly
// compatible with the MySQL protocol.
//
// Although TiDB is compatible with MySQL protocol, it is packaged as a separate driver component
// rather than reusing the mysql adapter directly. This design allows for future extensibility,
// such as implementing TiDB-specific features like distributed transactions or optimizations.
type Driver struct {
	*mysql.Driver
}

func init() {
	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"tidb"}
	)
	for _, driverName := range driverNames {
		if err = gdb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
}

// New creates and returns a driver that implements gdb.Driver, which supports operations for TiDB.
func New() gdb.Driver {
	mysqlDriver := mysql.New().(*mysql.Driver)
	return &Driver{
		Driver: mysqlDriver,
	}
}
