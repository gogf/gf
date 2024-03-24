// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmetric"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	instrumentName                              = "github.com/gogf/gf/v2/net/gclient.Client"
	tracingAttrHttpAddressRemote                = "http.address.remote"
	tracingAttrHttpAddressLocal                 = "http.address.local"
	tracingAttrHttpDnsStart                     = "http.dns.start"
	tracingAttrHttpDnsDone                      = "http.dns.done"
	tracingAttrHttpConnectStart                 = "http.connect.start"
	tracingAttrHttpConnectDone                  = "http.connect.done"
	tracingEventHttpRequest                     = "http.request"
	tracingEventHttpRequestHeaders              = "http.request.headers"
	tracingEventHttpRequestBaggage              = "http.request.baggage"
	tracingEventHttpRequestBody                 = "http.request.body"
	tracingEventHttpResponse                    = "http.response"
	tracingEventHttpResponseHeaders             = "http.response.headers"
	tracingEventHttpResponseBody                = "http.response.body"
	tracingMiddlewareHandled        gctx.StrKey = `MiddlewareClientTracingHandled`
)

// internalMiddlewareObservability is a client middleware that enables observability feature.
func internalMiddlewareObservability(c *Client, r *http.Request) (response *Response, err error) {
	var ctx = r.Context()
	// Mark this request is handled by server tracing middleware,
	// to avoid repeated handling by the same middleware.
	if ctx.Value(tracingMiddlewareHandled) != nil {
		return c.Next(r)
	}

	ctx = context.WithValue(ctx, tracingMiddlewareHandled, 1)
	tr := otel.GetTracerProvider().Tracer(
		instrumentName,
		trace.WithInstrumentationVersion(gf.VERSION),
	)
	ctx, span := tr.Start(ctx, r.URL.String(), trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(gtrace.CommonLabels()...)

	// Inject tracing content into http header.
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

	// Inject ClientTrace into context for http request.
	var (
		httpClientTracer       *httptrace.ClientTrace
		baseClientTracer       = newClientTracerNoop()
		isUsingDefaultProvider = gtrace.IsUsingDefaultProvider()
	)
	// Tracing.
	if !isUsingDefaultProvider {
		baseClientTracer = newClientTracerTracing(ctx, span, r)
	}
	// Metrics.
	if gmetric.IsEnabled() {
		baseClientTracer = newClientTracerMetrics(r, baseClientTracer)
	}
	httpClientTracer = newClientTracer(baseClientTracer)
	r = r.WithContext(
		httptrace.WithClientTrace(
			ctx, httpClientTracer,
		),
	)
	response, err = c.Next(r)

	// If it is now using default trace provider, it then does no complex tracing jobs.
	if isUsingDefaultProvider {
		return
	}

	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
	if response == nil || response.Response == nil {
		return
	}

	reqBodyContentBytes, _ := io.ReadAll(response.Body)
	response.Body = utils.NewReadCloser(reqBodyContentBytes, false)

	resBodyContent, err := gtrace.SafeContentForHttp(reqBodyContentBytes, response.Header)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`converting safe content failed: %s`, err.Error()))
	}

	span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(
		attribute.String(tracingEventHttpResponseHeaders, gconv.String(httputil.HeaderToMap(response.Header))),
		attribute.String(tracingEventHttpResponseBody, resBodyContent),
	))
	return
}
