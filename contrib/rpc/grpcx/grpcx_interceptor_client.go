// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/tracing"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryError handles the error types converting between grpc and gerror.
// Note that, the minus error code is only used locally which will not be sent to other side.
func (c modClient) UnaryError(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		grpcStatus, ok := status.FromError(err)
		if ok {
			if code := grpcStatus.Code(); code != 0 {
				return gerror.NewCode(gcode.New(int(code), "", nil), grpcStatus.Message())
			}
			return gerror.New(grpcStatus.Message())
		}
	}
	return err
}

// UnaryTracing is a unary interceptor for adding tracing feature for gRPC client using OpenTelemetry.
func (c modClient) UnaryTracing(
	ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return tracing.UnaryClientInterceptor(ctx, method, req, reply, cc, invoker, opts...)
}

// StreamTracing is a stream interceptor for adding tracing feature for gRPC client using OpenTelemetry.
func (c modClient) StreamTracing(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer,
	callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
	return tracing.StreamClientInterceptor(ctx, desc, cc, method, streamer, callOpts...)
}
