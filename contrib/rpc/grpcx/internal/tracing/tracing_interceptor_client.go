// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/grpcctx"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/utils"
	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/util/gconv"
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
	tracer := otel.GetTracerProvider().Tracer(
		tracingInstrumentGrpcClient,
		trace.WithInstrumentationVersion(gf.VERSION),
	)
	requestMetadata, _ := metadata.FromOutgoingContext(ctx)
	metadataCopy := requestMetadata.Copy()
	name, attr := spanInfo(method, cc.Target())
	var span trace.Span
	ctx, span = tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...),
	)
	defer span.End()

	Inject(ctx, metadataCopy)

	ctx = metadata.NewOutgoingContext(ctx, metadataCopy)

	// If it is now using default trace provider, it then does no complex tracing jobs.
	if gtrace.IsUsingDefaultProvider() {
		return invoker(ctx, method, req, reply, cc, callOpts...)
	}

	span.SetAttributes(gtrace.CommonLabels()...)

	span.AddEvent(tracingEventGrpcRequest, trace.WithAttributes(
		attribute.String(tracingEventGrpcRequestBaggage, gconv.String(gtrace.GetBaggageMap(ctx))),
		attribute.String(tracingEventGrpcMetadataOutgoing, gconv.String(grpcctx.Ctx{}.OutgoingMap(ctx))),
		attribute.String(
			tracingEventGrpcRequestMessage,
			utils.MarshalMessageToJsonStringForTracing(
				req, "Request", tracingMaxContentLogSize,
			),
		),
	))

	err := invoker(ctx, method, req, reply, cc, callOpts...)

	span.AddEvent(tracingEventGrpcResponse, trace.WithAttributes(
		attribute.String(
			tracingEventGrpcResponseMessage,
			utils.MarshalMessageToJsonStringForTracing(
				reply, "Response", tracingMaxContentLogSize,
			),
		),
	))

	if err != nil {
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		span.SetAttributes(statusCodeAttr(grpcCodes.OK))
	}
	return err
}

// StreamClientInterceptor returns a grpc.StreamClientInterceptor suitable
// for use in a grpc.Dial call.
func StreamClientInterceptor(
	ctx context.Context, desc *grpc.StreamDesc,
	cc *grpc.ClientConn, method string, streamer grpc.Streamer,
	callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
	tracer := otel.GetTracerProvider().Tracer(
		tracingInstrumentGrpcClient,
		trace.WithInstrumentationVersion(gf.VERSION),
	)
	requestMetadata, _ := metadata.FromOutgoingContext(ctx)
	metadataCopy := requestMetadata.Copy()
	name, attr := spanInfo(method, cc.Target())

	var span trace.Span
	ctx, span = tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attr...),
	)

	Inject(ctx, metadataCopy)
	ctx = metadata.NewOutgoingContext(ctx, metadataCopy)

	span.SetAttributes(gtrace.CommonLabels()...)

	s, err := streamer(ctx, desc, cc, method, callOpts...)
	stream := wrapClientStream(s, desc)
	go func() {
		if err == nil {
			err = <-stream.finished
		}

		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.SetAttributes(statusCodeAttr(s.Code()))
		} else {
			span.SetAttributes(statusCodeAttr(grpcCodes.OK))
		}

		span.End()
	}()

	return stream, err
}
