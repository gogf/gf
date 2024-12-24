// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package internal

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
)

// InitTracer initializes and registers `otlpgrpc` or `otlphttp` to global TracerProvider.
func InitTracer(opts ...trace.TracerProviderOption) (func(ctx context.Context), error) {
	tracerProvider := trace.NewTracerProvider(opts...)
	// Set the global propagator to traceContext (not set by default).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	return func(ctx context.Context) {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		// Shutdown waits for exported trace spans to be uploaded.
		if err := tracerProvider.Shutdown(ctx); err != nil {
			g.Log().Errorf(ctx, "Shutdown tracerProvider failed err:%+v", err)
		} else {
			g.Log().Debug(ctx, "Shutdown tracerProvider success")
		}
	}, nil
}

// GetLocalIP returns the IP address of the server.
func GetLocalIP() (string, error) {
	var intranetIPArray, err = gipv4.GetIntranetIpArray()
	if err != nil {
		return "", err
	}

	if len(intranetIPArray) == 0 {
		if intranetIPArray, err = gipv4.GetIpArray(); err != nil {
			return "", err
		}
	}
	var hostIP = "NoHostIpFound"
	if len(intranetIPArray) > 0 {
		hostIP = intranetIPArray[0]
	}
	return hostIP, nil
}
