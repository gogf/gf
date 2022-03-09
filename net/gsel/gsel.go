// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsel provides selector definition and implements.
package gsel

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
)

// Builder creates and returns selector in runtime.
type Builder interface {
	Build() Selector
}

// Selector for service balancer.
type Selector interface {
	// Pick selects and returns service.
	Pick(ctx context.Context) (node Node, done DoneFunc, err error)

	// Update updates services into Selector.
	Update(nodes []Node) error
}

// Node is node interface.
type Node interface {
	Service() *gsvc.Service
	Address() string
}

// DoneFunc is callback function when RPC invoke done.
type DoneFunc func(ctx context.Context, di DoneInfo)

// DoneInfo contains additional information for done.
type DoneInfo struct {
	// Err is the rpc error the RPC finished with. It could be nil.
	Err error

	// Trailer contains the metadata from the RPC's trailer, if present.
	Trailer DoneInfoMD

	// BytesSent indicates if any bytes have been sent to the server.
	BytesSent bool

	// BytesReceived indicates if any byte has been received from the server.
	BytesReceived bool

	// ServerLoad is the load received from server. It's usually sent as part of
	// trailing metadata.
	//
	// The only supported type now is *orca_v1.LoadReport.
	ServerLoad interface{}
}

// DoneInfoMD is a mapping from metadata keys to value array.
// Users should use the following two convenience functions New and Pairs to generate MD.
type DoneInfoMD interface {
	// Len returns the number of items in md.
	Len() int

	// Get obtains the values for a given key.
	//
	// k is converted to lowercase before searching in md.
	Get(k string) []string

	// Set sets the value of a given key with a slice of values.
	//
	// k is converted to lowercase before storing in md.
	Set(key string, values ...string)

	// Append adds the values to key k, not overwriting what was already stored at
	// that key.
	//
	// k is converted to lowercase before storing in md.
	Append(k string, values ...string)

	// Delete removes the values for a given key k which is converted to lowercase
	// before removing it from md.
	Delete(k string)
}
