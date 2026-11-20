// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// UnaryLogger is the default unary interceptor for logging purpose.
func (s *GrpcServer) UnaryLogger(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	var (
		start    = time.Now()
		res, err = handler(ctx, req)
		duration = time.Since(start)
	)
	s.handleAccessLog(ctx, err, duration, info, req, res)
	s.handleErrorLog(ctx, err, duration, info, req, res)
	return res, err
}

// handleAccessLog handles the access logging for server.
func (s *GrpcServer) handleAccessLog(
	ctx context.Context, err error, duration time.Duration, info *grpc.UnaryServerInfo, req, res any,
) {
	if !s.config.AccessLogEnabled {
		return
	}
	content := fmt.Sprintf(
		"%s, %.3fms, %+v, %+v",
		info.FullMethod, float64(duration)/1e6, req, res,
	)
	s.config.Logger.Stdout(s.config.LogStdout).File(s.config.AccessLogPattern).Path(s.config.LogPath).Print(ctx, content)
}

// handleErrorLog handles the error logging for server.
func (s *GrpcServer) handleErrorLog(
	ctx context.Context, err error, duration time.Duration, info *grpc.UnaryServerInfo, req, res any,
) {
	// It does nothing if error logging is custom disabled.
	if !s.config.ErrorLogEnabled || err == nil {
		return
	}
	var (
		code          = gerror.Code(err)
		codeDetail    = code.Detail()
		codeDetailStr string
		grpcCode      codes.Code
		grpcMessage   string
	)
	if grpcStatus, ok := status.FromError(err); ok {
		grpcCode = grpcStatus.Code()
		grpcMessage = grpcStatus.Message()
	}
	if codeDetail != nil {
		codeDetailStr = gstr.Replace(fmt.Sprintf(`%+v`, codeDetail), "\n", " ")
	}
	content := fmt.Sprintf(
		`%s, %.3fms, %d, "%s", %+v, %+v, %d, "%s", "%s"`,
		info.FullMethod, float64(duration)/1e6, grpcCode, grpcMessage,
		req, res, code.Code(), code.Message(), codeDetailStr,
	)
	if s.config.ErrorStack {
		if stack := gerror.Stack(err); stack != "" {
			content += "\nStack:\n" + stack
		} else {
			content += ", " + err.Error()
		}
	} else {
		content += ", " + err.Error()
	}
	s.config.Logger.Stack(false).
		Stdout(s.config.LogStdout).
		File(s.config.ErrorLogPattern).Path(s.config.LogPath).Error(ctx, content)
}

// UnaryTimeout is the unary interceptor that limits the handling duration of a request using
// `Timeout` and `MethodTimeouts` of the server configuration. It does nothing if no timeout applies.
//
// The handler runs in a separate goroutine so that the request returns with `DeadlineExceeded`
// on time even if the handler blocks and ignores the context.
func (s *GrpcServer) UnaryTimeout(
	ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	timeout, ok := s.config.MethodTimeouts[info.FullMethod]
	if !ok {
		timeout = s.config.Timeout
	}
	if timeout <= 0 {
		return handler(ctx, req)
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	type result struct {
		res      any
		err      error
		panicked any
	}
	// Buffered, so that the goroutine never leaks if the request has already timed out.
	resultChan := make(chan result, 1)
	go func() {
		var r result
		defer func() {
			if p := recover(); p != nil {
				r.panicked = p
			}
			resultChan <- r
		}()
		r.res, r.err = handler(ctx, req)
	}()

	select {
	case r := <-resultChan:
		if r.panicked != nil {
			panic(r.panicked)
		}
		return r.res, r.err
	case <-ctx.Done():
		// Prefer the real result if it is ready at the same moment.
		select {
		case r := <-resultChan:
			if r.panicked != nil {
				panic(r.panicked)
			}
			return r.res, r.err
		default:
		}
		return nil, status.FromContextError(ctx.Err()).Err()
	}
}
