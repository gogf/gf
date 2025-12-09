// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gaussdb implements gdb.Driver, which supports operations for database GaussDB.
package gaussdb

import (
	"github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Driver is the driver for GaussDB database.
//
// GaussDB is an enterprise-level distributed database developed by Huawei. GaussDB for MySQL is a cloud-native
// database that is fully compatible with MySQL protocol.
//
// Although GaussDB is compatible with MySQL protocol, it is packaged as a separate driver component
// rather than reusing the mysql adapter directly. This design allows for future extensibility,
// such as implementing GaussDB-specific features or optimizations for cloud-native scenarios.
type Driver struct {
	*mysql.Driver
}

func init() {
	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"gaussdb"}
	)
	for _, driverName := range driverNames {
		if err = gdb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
}

// New creates and returns a driver that implements gdb.Driver, which supports operations for GaussDB.
func New() gdb.Driver {
	mysqlDriver := mysql.New().(*mysql.Driver)
	return &Driver{
		Driver: mysqlDriver,
	}
}
