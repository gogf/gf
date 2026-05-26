// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpcx provides gapp Server adapter for GrpcServer.
package grpcx

import (
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gapp"
)

// adapterStartTimeout is the maximum duration to wait for a gRPC server
// to indicate it is listening after Start() is called.
const grpcAdapterStartTimeout = time.Second * 2

// adapterStartPollInterval is the interval between readiness polls.
const grpcAdapterStartPollInterval = time.Millisecond * 10

// GrpcServerAdapter wraps GrpcServer to implement the gapp.Server interface.
type GrpcServerAdapter struct {
	server *GrpcServer
}

// NewGappServerAdapter creates and returns a gapp.Server adapter for GrpcServer.
func NewGappServerAdapter(server *GrpcServer) gapp.Server {
	return &GrpcServerAdapter{server: server}
}

// Start starts the gRPC server in non-blocking way without registering signal handlers.
func (a *GrpcServerAdapter) Start() error {
	if err := a.server.StartManaged(); err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err, "grpc server start failed")
	}

	// Poll until the listener is established or timeout.
	deadline := time.Now().Add(grpcAdapterStartTimeout)
	for time.Now().Before(deadline) {
		if a.server.GetListenedPort() > 0 {
			return nil
		}
		time.Sleep(grpcAdapterStartPollInterval)
	}
	// serve() succeeded but we cannot confirm readiness; clean up the running server.
	a.server.StopForceful()
	return gerror.NewCode(gcode.CodeOperationFailed, "grpc server failed to start within timeout")
}

// Stop stops the gRPC server.
// When graceful is true, it waits for in-flight RPCs to complete.
// When graceful is false, it forcibly stops the server immediately.
// Note: GrpcServer.Stop/StopForceful do not return errors, so this
// method always returns nil.
func (a *GrpcServerAdapter) Stop(graceful bool) error {
	if graceful {
		a.server.Stop()
	} else {
		a.server.StopForceful()
	}
	return nil
}
