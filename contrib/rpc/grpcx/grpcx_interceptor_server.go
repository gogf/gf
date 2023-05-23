// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"context"
	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/tracing"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gutil"
)

// ChainUnary returns a ServerOption that specifies the chained interceptor
// for unary RPCs. The first interceptor will be the outermost,
// while the last interceptor will be the innermost wrapper around the real call.
// All unary interceptors added by this method will be chained.
func (s modServer) ChainUnary(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(interceptors...)
}

// ChainStream returns a ServerOption that specifies the chained interceptor
// for stream RPCs. The first interceptor will be the outermost,
// while the last interceptor will be the innermost wrapper around the real call.
// All stream interceptors added by this method will be chained.
func (s modServer) ChainStream(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.ChainStreamInterceptor(interceptors...)
}

// UnaryError is the default unary interceptor for error converting from custom error to grpc error.
func (s modServer) UnaryError(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		code := gerror.Code(err)
		if code.Code() != -1 {
			err = status.Error(codes.Code(code.Code()), err.Error())
		}
	}
	return res, err
}

// UnaryRecover is the first interceptor that keep server not down from panics.
func (s modServer) UnaryRecover(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (res interface{}, err error) {
	gutil.TryCatch(ctx, func(ctx2 context.Context) {
		res, err = handler(ctx, req)
	}, func(ctx context.Context, exception error) {
		err = gerror.WrapCode(gcode.New(int(codes.Internal), "", nil), err, "panic recovered")
	})
	return
}

// UnaryValidate Common validation unary interpreter.
func (s modServer) UnaryValidate(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	// It does nothing if there's no validation tag in the struct definition.
	if err := g.Validator().Data(req).Run(ctx); err != nil {
		return nil, gerror.NewCode(
			gcode.New(int(codes.InvalidArgument), "", nil),
			gerror.Current(err).Error(),
		)
	}
	return handler(ctx, req)
}

func (s modServer) UnaryAllowNilRes(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	res, err := handler(ctx, req)
	if g.IsNil(res) {
		res = proto.Message(nil)
	}
	return res, err
}

// UnaryTracing is a unary interceptor for adding tracing feature for gRPC server using OpenTelemetry.
func (s modServer) UnaryTracing(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	return tracing.UnaryServerInterceptor(ctx, req, info, handler)
}

// StreamTracing is a stream unary interceptor for adding tracing feature for gRPC server using OpenTelemetry.
func (s modServer) StreamTracing(
	srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) error {
	return tracing.StreamServerInterceptor(srv, ss, info, handler)
}
