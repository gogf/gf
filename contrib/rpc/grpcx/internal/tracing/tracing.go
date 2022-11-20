// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc/interceptor.go

package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

const (
	// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")
)

const (
	tracingMaxContentLogSize         = 256 * 1024 // Max log size for request and response body.
	tracingInstrumentGrpcClient      = "github.com/gogf/gf/contrib/rpc/grpcx/v2/krpc.GrpcClient"
	tracingInstrumentGrpcServer      = "github.com/gogf/gf/contrib/rpc/grpcx/v2/krpc.GrpcServer"
	tracingEventGrpcRequest          = "grpc.request"
	tracingEventGrpcRequestMessage   = "grpc.request.message"
	tracingEventGrpcRequestBaggage   = "grpc.request.baggage"
	tracingEventGrpcMetadataOutgoing = "grpc.metadata.outgoing"
	tracingEventGrpcMetadataIncoming = "grpc.metadata.incoming"
	tracingEventGrpcResponse         = "grpc.response"
	tracingEventGrpcResponseMessage  = "grpc.response.message"
)

type metadataSupplier struct {
	metadata metadata.MD
}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	var (
		index = 0
		keys  = make([]string, s.metadata.Len())
	)
	for k := range s.metadata {
		keys[index] = k
		index++
	}
	return keys
}

// Inject injects correlation context and span context into the gRPC
// metadata object. This function is meant to be used on outgoing
// requests.
func Inject(ctx context.Context, metadata metadata.MD) {
	otel.GetTextMapPropagator().Inject(ctx, &metadataSupplier{
		metadata: metadata,
	})
}

// Extract returns the correlation context and span context that
// another service encoded in the gRPC metadata object with Inject.
// This function is meant to be used on incoming requests.
func Extract(ctx context.Context, metadata metadata.MD) (baggage.Baggage, trace.SpanContext) {
	ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataSupplier{
		metadata: metadata,
	})
	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}
