// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package jaeger provides OpenTelemetry provider using jaeger.
//
// Deprecated:
// Please use `github.com/gogf/gf/contrib/trace/otlphttp/v2` instead.
// As dependent package `go.opentelemetry.io/otel/exporters/jaeger` is deprecated in otel,
// this package will be never be maintained from goframe v2.6.0.
package jaeger

import (
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/text/gregex"
)

const (
	tracerHostnameTagKey = "hostname"
)

// Init initializes and registers jaeger to global TracerProvider.
//
// The output parameter `flush` is used for waiting exported trace spans to be uploaded,
// which is useful if your program is ending, and you do not want to lose recent spans.
func Init(serviceName, endpoint string) (*trace.TracerProvider, error) {
	// Create the Jaeger exporter
	var endpointOption jaeger.EndpointOption
	if strings.HasPrefix(endpoint, "http") {
		// HTTP.
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
	} else {
		// UDP.
		match, err := gregex.MatchString(`(.+):(\d+)`, endpoint)
		if err != nil {
			return nil, err
		}
		if len(match) < 3 {
			return nil, gerror.NewCodef(
				gcode.CodeInvalidParameter, `invalid endpoint "%s"`, endpoint,
			)
		}
		var (
			host = match[1]
			port = match[2]
		)
		endpointOption = jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(host), jaeger.WithAgentPort(port),
		)
	}
	// Try retrieving host ip for tracing info.
	var (
		hostIp               = "NoHostIpFound"
		intranetIPArray, err = gipv4.GetIntranetIpArray()
	)
	if err != nil {
		return nil, err
	}
	if len(intranetIPArray) == 0 {
		if intranetIPArray, err = gipv4.GetIpArray(); err != nil {
			return nil, err
		}
	}
	if len(intranetIPArray) > 0 {
		hostIp = intranetIPArray[0]
	}

	exp, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.HostNameKey.String(hostIp),
			attribute.String(tracerHostnameTagKey, hostIp),
		)),
	)
	// Register our TracerProvider as the global, so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	return tp, nil
}
