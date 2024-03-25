// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	instrumentName                              = "github.com/gogf/gf/v2/net/ghttp.Server"
	tracingEventHttpRequest                     = "http.request"
	tracingEventHttpRequestHeaders              = "http.request.headers"
	tracingEventHttpRequestBaggage              = "http.request.baggage"
	tracingEventHttpRequestBody                 = "http.request.body"
	tracingEventHttpResponse                    = "http.response"
	tracingEventHttpResponseHeaders             = "http.response.headers"
	tracingEventHttpResponseBody                = "http.response.body"
	tracingEventHttpRequestUrl                  = "http.request.url"
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

	// Request content logging.
	reqBodyContentBytes, err := io.ReadAll(r.Body)
	if err != nil {
		r.SetError(gerror.Wrap(err, `read request body failed`))
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
		return
	}
	r.Body = utils.NewReadCloser(reqBodyContentBytes, false)
	reqBodyContent, err := gtrace.SafeContentForHttp(reqBodyContentBytes, r.Header)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`converting safe content failed: %s`, err.Error()))
	}

	span.AddEvent(tracingEventHttpRequest, trace.WithAttributes(
		attribute.String(tracingEventHttpRequestUrl, r.URL.String()),
		attribute.String(tracingEventHttpRequestHeaders, gconv.String(httputil.HeaderToMap(r.Header))),
		attribute.String(tracingEventHttpRequestBaggage, gtrace.GetBaggageMap(ctx).String()),
		attribute.String(tracingEventHttpRequestBody, reqBodyContent),
	))

	// Continue executing.
	r.Middleware.Next()

	// parse after set route as span name
	if handler := r.GetServeHandler(); handler != nil && handler.Handler.Router != nil {
		span.SetName(handler.Handler.Router.Uri)
	}

	// Error logging.
	if err = r.GetError(); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}

	// Response content logging.
	resBodyContent, err := gtrace.SafeContentForHttp(r.Response.Buffer(), r.Response.Header())
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`converting safe content failed: %s`, err.Error()))
	}

	span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(
		attribute.String(tracingEventHttpResponseHeaders, gconv.String(httputil.HeaderToMap(r.Response.Header()))),
		attribute.String(tracingEventHttpResponseBody, resBodyContent),
	))
}
