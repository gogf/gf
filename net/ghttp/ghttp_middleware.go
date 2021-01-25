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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
)

const (
	tracingMaxContentLogSize = 512 * 1024 // Max log size for request and response body.
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
	ctx, span := tr.Start(ctx, r.URL.String())
	defer span.End()

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
			"[Request Body Too Large For Logging, Max: %d bytes]",
			tracingMaxContentLogSize,
		)
	}
	span.AddEvent("http.request", trace.WithAttributes(
		label.Any(`http.request.headers`, httputil.HeaderToMap(r.Header)),
		label.String(`http.request.body`, reqBodyContent),
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
		resBodyContent = fmt.Sprintf("[Response Body Too Large For Logging, Max: %d bytes]", tracingMaxContentLogSize)
	}
	span.AddEvent("http.response", trace.WithAttributes(
		label.Any(`http.response.headers`, httputil.HeaderToMap(r.Response.Header())),
		label.String(`http.response.body`, resBodyContent),
	))
	return
}
