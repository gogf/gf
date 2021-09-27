// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/container/gvar"
)

// Adapter is an interface for universal redis operations.
type Adapter interface {
	// Conn retrieves and returns a connection object for continuous operations.
	// Note that you should call Close function manually if you do not use this connection any further.
	Conn(ctx context.Context) (conn Conn, err error)

	// Close closes current redis client, closes its connection pool and releases all its related resources.
	Close(ctx context.Context) (err error)
}

// Conn is an interface of a connection from universal redis client.
type Conn interface {
	// Do sends a command to the server and returns the received reply.
	// It uses json.Marshal for struct/slice/map type values before committing them to redis.
	Do(ctx context.Context, command string, args ...interface{}) (result *gvar.Var, err error)

	// Receive receives a single reply as gvar.Var from the Redis server.
	Receive(ctx context.Context) (result *gvar.Var, err error)

	// Close puts the connection back to connection pool.
	Close(ctx context.Context) (err error)
}
