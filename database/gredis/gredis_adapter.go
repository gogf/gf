// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/container/gvar"
	"time"
)

type Option struct {
	ReadTimeout time.Duration
}

type Adapter interface {
	Conn(ctx context.Context) (conn Conn, err error)
	Stats(ctx context.Context) (stats Stats, err error)
	Close(ctx context.Context) (err error)
}

type Conn interface {
	Do(ctx context.Context, command string, args []interface{}, option *Option) (result *gvar.Var, err error)
	Receive(ctx context.Context, option *Option) (result *gvar.Var, err error)
	Close(ctx context.Context) (err error)
}

type Stats interface {
	// ActiveCount is the number of connections in the pool. The count includes
	// idle connections and connections in use.
	ActiveCount() int64

	// IdleCount is the number of idle connections in the pool.
	IdleCount() int64

	// WaitCount is the total number of connections waited for.
	// This value is currently not guaranteed to be 100% accurate.
	WaitCount() int64
}
