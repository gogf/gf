// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/net/ghttp/internal/httputil"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracingInstrumentName           = "github.com/gogf/gf/v2/net/ghttp.Client"
	tracingAttrHttpAddressRemote    = "http.address.remote"
	tracingAttrHttpAddressLocal     = "http.address.local"
	tracingAttrHttpDnsStart         = "http.dns.start"
	tracingAttrHttpDnsDone          = "http.dns.done"
	tracingAttrHttpConnectStart     = "http.connect.start"
	tracingAttrHttpConnectDone      = "http.connect.done"
	tracingEventHttpRequest         = "http.request"
	tracingEventHttpRequestHeaders  = "http.request.headers"
	tracingEventHttpRequestBaggage  = "http.request.baggage"
	tracingEventHttpRequestBody     = "http.request.body"
	tracingEventHttpResponse        = "http.response"
	tracingEventHttpResponseHeaders = "http.response.headers"
	tracingEventHttpResponseBody    = "http.response.body"
)

// MiddlewareTracing is a client middleware that enables tracing feature using standards of OpenTelemetry.
func MiddlewareTracing(c *Client, r *http.Request) (response *Response, err error) {
	tr := otel.GetTracerProvider().Tracer(tracingInstrumentName, trace.WithInstrumentationVersion(gf.VERSION))
	ctx, span := tr.Start(r.Context(), r.URL.String(), trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(gtrace.CommonLabels()...)

	// Inject tracing content into http header.
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

	// Continue client handler executing.
	response, err = c.Next(
		r.WithContext(
			httptrace.WithClientTrace(
				ctx, newClientTrace(ctx, span, r),
			),
		),
	)
	if err != nil {
		span.SetStatus(codes.Error, fmt.Sprintf(`%+v`, err))
	}
	if response == nil || response.Response == nil {
		return
	}

	reqBodyContentBytes, _ := ioutil.ReadAll(response.Body)
	response.Body = utils.NewReadCloser(reqBodyContentBytes, false)

	span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(
		attribute.String(tracingEventHttpResponseHeaders, gconv.String(httputil.HeaderToMap(response.Header))),
		attribute.String(tracingEventHttpResponseBody, gstr.StrLimit(
			string(reqBodyContentBytes),
			gtrace.MaxContentLogSize(),
			"...",
		)),
	))
	return
}
