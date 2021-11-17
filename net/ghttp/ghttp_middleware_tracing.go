// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/net/ghttp/internal/client"
	"github.com/gogf/gf/v2/net/ghttp/internal/httputil"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	tracingInstrumentName           = "github.com/gogf/gf/v2/net/ghttp.Server"
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
	var (
		tr        = otel.GetTracerProvider().Tracer(tracingInstrumentName, trace.WithInstrumentationVersion(gf.VERSION))
		ctx, span = tr.Start(
			otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header)),
			r.URL.String(),
			trace.WithSpanKind(trace.SpanKindServer),
		)
	)
	defer span.End()

	span.SetAttributes(gtrace.CommonLabels()...)

	// Inject tracing context.
	r.SetCtx(ctx)

	// Request content logging.
	reqBodyContentBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = utils.NewReadCloser(reqBodyContentBytes, false)

	span.AddEvent(tracingEventHttpRequest, trace.WithAttributes(
		attribute.String(tracingEventHttpRequestHeaders, gconv.String(httputil.HeaderToMap(r.Header))),
		attribute.String(tracingEventHttpRequestBaggage, gtrace.GetBaggageMap(ctx).String()),
		attribute.String(tracingEventHttpRequestBody, gstr.StrLimit(
			string(reqBodyContentBytes),
			gtrace.MaxContentLogSize(),
			"...",
		)),
	))

	// Continue executing.
	r.Middleware.Next()

	// Error logging.
	if err := r.GetError(); err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
	// Response content logging.
	var resBodyContent = gstr.StrLimit(r.Response.BufferString(), gtrace.MaxContentLogSize(), "...")

	span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(
		attribute.String(tracingEventHttpResponseHeaders, gconv.String(httputil.HeaderToMap(r.Response.Header()))),
		attribute.String(tracingEventHttpResponseBody, resBodyContent),
	))
	return
}
