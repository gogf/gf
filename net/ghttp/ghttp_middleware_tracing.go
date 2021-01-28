// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"github.com/gogf/gf"
	"github.com/gogf/gf/internal/utils"
	"github.com/gogf/gf/net/ghttp/internal/client"
	"github.com/gogf/gf/net/ghttp/internal/httputil"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
)

const (
	tracingMaxContentLogSize        = 512 * 1024 // Max log size for request and response body.
	tracingEventHttpRequest         = "http.request"
	tracingEventHttpRequestHeaders  = "http.request.headers"
	tracingEventHttpRequestBaggage  = "http.request.baggage"
	tracingEventHttpRequestBody     = "http.request.body"
	tracingEventHttpResponse        = "http.response"
	tracingEventHttpResponseHeaders = "http.response.headers"
	tracingEventHttpResponseBody    = "http.response.body"
)

// MiddlewareClientTracing is a client middleware that enables tracing feature using standards of OpenTelemetry.
func MiddlewareClientTracing(c *Client, r *http.Request) (*ClientResponse, error) {
	return client.MiddlewareTracing(c, r)
}

// MiddlewareServerTracing is a serer middleware that enables tracing feature using standards of OpenTelemetry.
func MiddlewareServerTracing(r *Request) {
	tr := otel.GetTracerProvider().Tracer(
		"github.com/gogf/gf/net/ghttp.Server",
		trace.WithInstrumentationVersion(fmt.Sprintf(`%s`, gf.VERSION)),
	)
	// Tracing content parsing, start root span.
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	ctx := propagator.Extract(r.Context(), r.Header)
	ctx, span := tr.Start(ctx, r.URL.String(), trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	span.SetAttributes(gtrace.CommonLabels()...)

	// Inject tracing context.
	r.SetCtx(ctx)

	// Request content logging.
	var reqBodyContent string
	if r.ContentLength <= tracingMaxContentLogSize {
		reqBodyContentBytes, _ := ioutil.ReadAll(r.Body)
		r.Body = utils.NewReadCloser(reqBodyContentBytes, false)
		reqBodyContent = string(reqBodyContentBytes)
	} else {
		reqBodyContent = fmt.Sprintf(
			"[Request Body Too Large For Tracing, Max: %d bytes]",
			tracingMaxContentLogSize,
		)
	}
	span.AddEvent(tracingEventHttpRequest, trace.WithAttributes(
		label.Any(tracingEventHttpRequestHeaders, httputil.HeaderToMap(r.Header)),
		label.Any(tracingEventHttpRequestBaggage, gtrace.GetBaggageMap(ctx).Map()),
		label.String(tracingEventHttpRequestBody, reqBodyContent),
	))

	// Continue executing.
	r.Middleware.Next()

	// Error logging.
	if err := r.GetError(); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
	// Response content logging.
	var resBodyContent string
	if r.Response.BufferLength() <= tracingMaxContentLogSize {
		resBodyContent = r.Response.BufferString()
	} else {
		resBodyContent = fmt.Sprintf(
			"[Response Body Too Large For Tracing, Max: %d bytes]",
			tracingMaxContentLogSize,
		)
	}
	span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(
		label.Any(tracingEventHttpResponseHeaders, httputil.HeaderToMap(r.Response.Header())),
		label.String(tracingEventHttpResponseBody, resBodyContent),
	))
	return
}
