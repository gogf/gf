// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	mysqlDriverName = "hellosql"
	quoteChar       = "`"
)

func init() {
	var (
		err       error
		driverObj = &DriverMysql{
			Driver: mysql.Driver{},
		}
	)
	if err = gdb.Register(mysqlDriverName, driverObj); err != nil {
		panic(err)
	}
}

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverMysql) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &DriverMysql{
		Driver: mysql.Driver{Core: core},
	}, nil
}

// GetChars returns the security char for this type of database.
func (d *DriverMysql) GetChars() (charLeft string, charRight string) {
	return quoteChar, quoteChar
}

func (d *DriverMysql) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	fmt.Println("DriverMysql.Open")
	fmt.Println("config.Pass(encode):" + config.Pass)
	// Decrypt the password if it is encrypted.
	config.Pass, err = gbase64.DecodeToString(config.Pass)
	if err != nil {
		return nil, err
	}
	fmt.Println("config.Pass(decode):" + config.Pass)
	return d.Driver.Open(config)
}

func (d *DriverMysql) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	return d.Driver.Tables(ctx, schema...)
}

func (d *DriverMysql) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	return d.Driver.TableFields(ctx, table, schema...)
}

// DriverMysql is the driver for mysql database.
type DriverMysql struct {
	mysql.Driver
}

func main() {
	list, err := g.DB().Tables(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}
