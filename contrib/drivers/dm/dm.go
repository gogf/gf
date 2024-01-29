// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dm implements gdb.Driver, which supports operations for database DM.
package dm

import (
	"context"
	"time"

	_ "gitee.com/chunanyong/dm"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type Driver struct {
	*gdb.Core
}

const (
	quoteChar = `"`
)

func init() {
	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"dm"}
	)
	for _, driverName := range driverNames {
		if err = gdb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for dm.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for dm.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return quoteChar, quoteChar
}

// ConvertValueForField converts value to the type of the record field.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) {
	switch itemValue := fieldValue.(type) {
	// dm does not support time.Time, it so here converts it to time string that it supports.
	case time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue.IsZero() {
			return nil, nil
		}
		return gtime.New(itemValue).String(), nil

	// dm does not support time.Time, it so here converts it to time string that it supports.
	case *time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue == nil || itemValue.IsZero() {
			return nil, nil
		}
		return gtime.New(itemValue).String(), nil
	}

	return fieldValue, nil
}

// TODO I originally wanted to only convert keywords in select
// 但是我发现 DoQuery 中会对 sql 会对 " " 达梦的安全字符 进行 / 转义，最后还是导致达梦无法正常解析
// However, I found that DoQuery() will perform / escape on sql with " " Dameng's safe characters, which ultimately caused Dameng to be unable to parse normally.
// But processing in DoFilter() is OK
// func (d *Driver) DoQuery(ctx context.Context, link gdb.Link, sql string, args ...interface{}) (gdb.Result, error) {
// 	l, r := d.GetChars()
// 	new := gstr.ReplaceI(sql, "INDEX", l+"INDEX"+r)
// 	g.Dump("new:", new)
// 	return d.Core.DoQuery(
// 		ctx,
// 		link,
// 		new,
// 		args,
// 	)
// }
