// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
)

// Adapter is an interface for universal redis operations.
type Adapter interface {
	AdapterGroup

	// Do send a command to the server and returns the received reply.
	// It uses json.Marshal for struct/slice/map type values before committing them to redis.
	Do(ctx context.Context, command string, args ...interface{}) (*gvar.Var, error)

	// Conn retrieves and returns a connection object for continuous operations.
	// Note that you should call Close function manually if you do not use this connection any further.
	Conn(ctx context.Context) (conn Conn, err error)

	// Close closes current redis client, closes its connection pool and releases all its related resources.
	Close(ctx context.Context) (err error)
}

// Conn is an interface of a connection from universal redis client.
type Conn interface {
	ConnCommand

	// Do send a command to the server and returns the received reply.
	// It uses json.Marshal for struct/slice/map type values before committing them to redis.
	Do(ctx context.Context, command string, args ...interface{}) (result *gvar.Var, err error)

	// Close puts the connection back to connection pool.
	Close(ctx context.Context) (err error)
}

// AdapterGroup is an interface managing group operations for redis.
type AdapterGroup interface {
	GroupGeneric() IGroupGeneric
	GroupHash() IGroupHash
	GroupList() IGroupList
	GroupPubSub() IGroupPubSub
	GroupScript() IGroupScript
	GroupSet() IGroupSet
	GroupSortedSet() IGroupSortedSet
	GroupString() IGroupString
}

// ConnCommand is an interface managing some operations bound to certain connection.
type ConnCommand interface {
	// Subscribe subscribes the client to the specified channels.
	// https://redis.io/commands/subscribe/
	Subscribe(ctx context.Context, channel string, channels ...string) ([]*Subscription, error)

	// PSubscribe subscribes the client to the given patterns.
	//
	// Supported glob-style patterns:
	// - h?llo subscribes to hello, hallo and hxllo
	// - h*llo subscribes to hllo and heeeello
	// - h[ae]llo subscribes to hello and hallo, but not hillo
	//
	// Use \ to escape special characters if you want to match them verbatim.
	//
	// https://redis.io/commands/psubscribe/
	PSubscribe(ctx context.Context, pattern string, patterns ...string) ([]*Subscription, error)

	// ReceiveMessage receives a single message of subscription from the Redis server.
	ReceiveMessage(ctx context.Context) (*Message, error)

	// Receive receives a single reply as gvar.Var from the Redis server.
	Receive(ctx context.Context) (result *gvar.Var, err error)
}
