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

// Selector for service balancer.
type Selector interface {
	// Pick selects and returns service.
	Pick(ctx context.Context) (node Node, done DoneFunc, err error)

	// Update updates services into Selector.
	Update(nodes []Node)
}

// Node is node interface.
type Node interface {
	Service() *gsvc.Service
}

// DoneFunc is callback function when RPC invoke done.
type DoneFunc func(ctx context.Context, di DoneInfo)

// DoneInfo contains additional information for done.
type DoneInfo struct {
	// Err is the rpc error the RPC finished with. It could be nil.
	Err error

	// Trailer contains the metadata from the RPC's trailer, if present.
	Trailer MD

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
