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
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
)

const (
	tracingMaxContentLogSize = 512 * 1024 // Max log size for request and response body.
)

// MiddlewareTracing is a client middleware that enables tracing feature using standards of OpenTelemetry.
func MiddlewareTracing(c *Client, r *http.Request) (response *Response, err error) {
	tr := otel.GetTracerProvider().Tracer(
		"github.com/gogf/gf/net/ghttp.Client",
		trace.WithInstrumentationVersion(fmt.Sprintf(`%s`, gf.VERSION)),
	)
	ctx, span := tr.Start(r.Context(), r.URL.String())
	defer span.End()

	span.SetAttributes(gtrace.CommonLabels()...)

	// Inject tracing content into http header.
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	propagator.Inject(ctx, r.Header)

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
		resBodyContent = fmt.Sprintf("[Response Body Too Large For Logging, Max: %d bytes]", tracingMaxContentLogSize)
	}

	span.AddEvent("http.response", trace.WithAttributes(
		label.Any(`http.response.headers`, httputil.HeaderToMap(response.Header)),
		label.String(`http.response.body`, resBodyContent),
	))
	return
}
