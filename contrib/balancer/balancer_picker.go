// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package balancer

import (
	"github.com/gogf/gf/v2/net/gsel"
	"google.golang.org/grpc/balancer"
)

// Picker implements grpc balancer.Picker,
// which is used by gRPC to pick a SubConn to send an RPC.
// Balancer is expected to generate a new picker from its snapshot every time its
// internal state has changed.
//
// The pickers used by gRPC can be updated by ClientConn.UpdateState().
type Picker struct {
	selector gsel.Selector
}

// Pick returns the connection to use for this RPC and related information.
//
// Pick should not block.  If the balancer needs to do I/O or any blocking
// or time-consuming work to service this call, it should return
// ErrNoSubConnAvailable, and the Pick call will be repeated by gRPC when
// the Picker is updated (using ClientConn.UpdateState).
//
// If an error is returned:
//
// - If the error is ErrNoSubConnAvailable, gRPC will block until a new
//   Picker is provided by the balancer (using ClientConn.UpdateState).
//
// - If the error is a status error (implemented by the grpc/status
//   package), gRPC will terminate the RPC with the code and message
//   provided.
//
// - For all other errors, wait for ready RPCs will wait, but non-wait for
//   ready RPCs will be terminated with this error's Error() string and
//   status code Unavailable.
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	node, done, err := p.selector.Pick(info.Ctx)
	if err != nil {
		return balancer.PickResult{}, err
	}
	return balancer.PickResult{
		SubConn: node.(*Node).conn,
		Done: func(di balancer.DoneInfo) {
			if done == nil {
				return
			}
			done(info.Ctx, gsel.DoneInfo{
				Err:           di.Err,
				Trailer:       di.Trailer,
				BytesSent:     di.BytesSent,
				BytesReceived: di.BytesReceived,
				ServerLoad:    di.ServerLoad,
			})
		},
	}, nil
}
