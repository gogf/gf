// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package clickhouse implements gdb.Driver, which supports operations for database ClickHouse.
package clickhouse

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

// Driver is the driver for postgresql database.
type Driver struct {
	*gdb.Core
}

var (
	errUnsupportedInsertIgnore = errors.New("unsupported method:InsertIgnore")
	errUnsupportedInsertGetId  = errors.New("unsupported method:InsertGetId")
	errUnsupportedReplace      = errors.New("unsupported method:Replace")
	errUnsupportedBegin        = errors.New("unsupported method:Begin")
	errUnsupportedTransaction  = errors.New("unsupported method:Transaction")
)

const (
	updateFilterPattern              = `(?i)UPDATE[\s]+?(\w+[\.]?\w+)[\s]+?SET`
	deleteFilterPattern              = `(?i)DELETE[\s]+?FROM[\s]+?(\w+[\.]?\w+)`
	filterTypePattern                = `(?i)^UPDATE|DELETE`
	replaceSchemaPattern             = `@(.+?)/([\w\.\-]+)+`
	needParsedSqlInCtx   gctx.StrKey = "NeedParsedSql"
	driverName                       = "clickhouse"
)

func init() {
	if err := gdb.Register(`clickhouse`, New()); err != nil {
		panic(err)
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for clickhouse.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for clickhouse.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

func (d *Driver) injectNeedParsedSql(ctx context.Context) context.Context {
	if ctx.Value(needParsedSqlInCtx) != nil {
		return ctx
	}
	return context.WithValue(ctx, needParsedSqlInCtx, true)
}
