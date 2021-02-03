// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package client

import (
	"fmt"
	"github.com/gogf/gf"
	"github.com/gogf/gf/internal/utils"
	"github.com/gogf/gf/net/ghttp/internal/httputil"
	"github.com/gogf/gf/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
)

const (
	tracingMaxContentLogSize        = 512 * 1024 // Max log size for request and response body.
	tracingInstrumentName           = "github.com/gogf/gf/net/ghttp.Client"
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
	otel.GetTextMapPropagator().Inject(ctx, r.Header)

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
	var resBodyContent string
	if response.ContentLength <= tracingMaxContentLogSize {
		reqBodyContentBytes, _ := ioutil.ReadAll(response.Body)
		resBodyContent = string(reqBodyContentBytes)
		response.Body = utils.NewReadCloser(reqBodyContentBytes, false)
	} else {
		resBodyContent = fmt.Sprintf(
			"[Response Body Too Large For Tracing, Max: %d bytes]",
			tracingMaxContentLogSize,
		)
	}

	span.AddEvent(tracingEventHttpResponse, trace.WithAttributes(
		label.Any(tracingEventHttpResponseHeaders, httputil.HeaderToMap(response.Header)),
		label.String(tracingEventHttpResponseBody, resBodyContent),
	))
	return
}
