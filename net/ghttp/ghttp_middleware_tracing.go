// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	instrumentName                              = "github.com/gogf/gf/v2/net/ghttp.Server"
	tracingEventHttpRequest                     = "http.request"
	tracingEventHttpRequestHeaders              = "http.request.headers"
	tracingEventHttpRequestBaggage              = "http.request.baggage"
	tracingEventHttpRequestParams               = "http.request.params"
	tracingEventHttpResponse                    = "http.response"
	tracingEventHttpResponseHeaders             = "http.response.headers"
	tracingEventHttpResponseBody                = "http.response.body"
	tracingEventHttpRequestUrl                  = "http.request.url"
	tracingEventHttpMethod                      = "http.method"
	tracingMiddlewareHandled        gctx.StrKey = `MiddlewareServerTracingHandled`
)

// internalMiddlewareServerTracing is a serer middleware that enables tracing feature using standards of OpenTelemetry.
func internalMiddlewareServerTracing(r *Request) {
	var (
		ctx = r.Context()
	)
	// Mark this request is handled by server tracing middleware,
	// to avoid repeated handling by the same middleware.
	if ctx.Value(tracingMiddlewareHandled) != nil {
		r.Middleware.Next()
		return
	}

	ctx = context.WithValue(ctx, tracingMiddlewareHandled, 1)
	var (
		span trace.Span
		tr   = otel.GetTracerProvider().Tracer(
			instrumentName,
			trace.WithInstrumentationVersion(gf.VERSION),
		)
	)
	ctx, span = tr.Start(
		otel.GetTextMapPropagator().Extract(
			ctx,
			propagation.HeaderCarrier(r.Header),
		),
		r.URL.Path,
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	span.SetAttributes(gtrace.CommonLabels()...)

	// Inject tracing context.
	r.SetCtx(ctx)

	// If it is now using a default trace provider, it then does no complex tracing jobs.
	if gtrace.IsUsingDefaultProvider() {
		r.Middleware.Next()
		return
	}

	// Basic trace attributes for all requests
	traceAttrs := []attribute.KeyValue{
		attribute.String(tracingEventHttpRequestUrl, r.URL.String()),
		attribute.String(tracingEventHttpMethod, r.Method),
		attribute.String(tracingEventHttpRequestHeaders, gconv.String(httputil.HeaderToMap(r.Header))),
		attribute.String(tracingEventHttpRequestBaggage, gtrace.GetBaggageMap(ctx).String()),
	}

	// Add request parameters if configured
	if r.Server != nil && r.Server.config.IsOtelTraceRequestEnabled() {
		// Get all request parameters (query + form + body)
		requestParams := make(map[string]any)

		// Query parameters
		for k, v := range r.URL.Query() {
			requestParams[k] = v
		}

		// Form parameters
		if r.ContentLength > 0 && gtrace.MaxContentLogSize() > 0 {
			contentType := r.Header.Get("Content-Type")
			if gstr.Contains(contentType, "application/x-www-form-urlencoded") || gstr.Contains(contentType, "multipart/form-data") {
				// Use GetFormMap() instead of ParseForm() to get form data
				formData := r.GetFormMap()
				for k, v := range formData {
					requestParams[k] = v
				}
			}
		}

		if len(requestParams) > 0 {
			traceAttrs = append(traceAttrs,
				attribute.String(tracingEventHttpRequestParams, gconv.String(requestParams)),
			)
		}
	}

	span.AddEvent(tracingEventHttpRequest, trace.WithAttributes(traceAttrs...))

	// Continue executing.
	r.Middleware.Next()

	// parse after set route as span name
	if handler := r.GetServeHandler(); handler != nil && handler.Handler.Router != nil {
		span.SetName(handler.Handler.Router.Uri)
	}

	// Error logging.
	if err := r.GetError(); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}

	// Response tracing attributes
	responseAttrs := []attribute.KeyValue{
		attribute.String(
			tracingEventHttpResponseHeaders,
			gconv.String(httputil.HeaderToMap(r.Response.Header())),
		),
	}

	// Add response body if configured
	if r.Server != nil && r.Server.config.IsOtelTraceResponseEnabled() {
		if r.Response.BufferLength() > 0 {
			responseBody := r.Response.BufferString()
			// Limit response body size for tracing to avoid memory issues
			if len(responseBody) > gtrace.MaxContentLogSize() {
				responseBody = responseBody[:gtrace.MaxContentLogSize()] + "...[truncated]"
			}
			responseAttrs = append(responseAttrs,
				attribute.String(tracingEventHttpResponseBody, responseBody),
			)
		}
	}

	span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(responseAttrs...))
}
