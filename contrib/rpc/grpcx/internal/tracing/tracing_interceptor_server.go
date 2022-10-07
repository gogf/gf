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
	"go.opentelemetry.io/otel/baggage"
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

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	tracer := otel.GetTracerProvider().Tracer(
		tracingInstrumentGrpcServer,
		trace.WithInstrumentationVersion(gf.VERSION),
	)
	requestMetadata, _ := metadata.FromIncomingContext(ctx)
	metadataCopy := requestMetadata.Copy()

	entries, spanCtx := Extract(ctx, metadataCopy)
	ctx = baggage.ContextWithBaggage(ctx, entries)
	ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
	name, attr := spanInfo(info.FullMethod, peerFromCtx(ctx))
	ctx, span := tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attr...),
	)
	defer span.End()

	// If it is now using default trace provider, it then does no complex tracing jobs.
	if gtrace.IsUsingDefaultProvider() {
		return handler(ctx, req)
	}

	span.SetAttributes(gtrace.CommonLabels()...)

	span.AddEvent(tracingEventGrpcRequest, trace.WithAttributes(
		attribute.String(tracingEventGrpcRequestBaggage, gconv.String(gtrace.GetBaggageMap(ctx))),
		attribute.String(tracingEventGrpcMetadataIncoming, gconv.String(grpcctx.Ctx{}.IncomingMap(ctx))),
		attribute.String(
			tracingEventGrpcRequestMessage,
			utils.MarshalMessageToJsonStringForTracing(
				req, "Request", tracingMaxContentLogSize,
			),
		),
	))

	res, err := handler(ctx, req)

	span.AddEvent(tracingEventGrpcResponse, trace.WithAttributes(
		attribute.String(
			tracingEventGrpcResponseMessage,
			utils.MarshalMessageToJsonStringForTracing(
				res, "Response", tracingMaxContentLogSize,
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

	return res, err
}

// StreamServerInterceptor returns a grpc.StreamServerInterceptor suitable
// for use in a grpc.NewServer call.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	tracer := otel.GetTracerProvider().Tracer(
		tracingInstrumentGrpcServer,
		trace.WithInstrumentationVersion(gf.VERSION),
	)

	ctx := ss.Context()
	requestMetadata, _ := metadata.FromIncomingContext(ctx)
	metadataCopy := requestMetadata.Copy()
	entries, spanCtx := Extract(ctx, metadataCopy)
	ctx = baggage.ContextWithBaggage(ctx, entries)
	ctx = trace.ContextWithRemoteSpanContext(ctx, spanCtx)
	name, attr := spanInfo(info.FullMethod, peerFromCtx(ctx))
	ctx, span := tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attr...),
	)
	defer span.End()

	span.SetAttributes(gtrace.CommonLabels()...)

	err := handler(srv, wrapServerStream(ctx, ss))

	if err != nil {
		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		span.SetAttributes(statusCodeAttr(s.Code()))
	} else {
		span.SetAttributes(statusCodeAttr(grpcCodes.OK))
	}

	return err
}
