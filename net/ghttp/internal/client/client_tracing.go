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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
)

const (
	maxContentLogSize = 512 * 1024 // Max log size for request and response body.
)

// MiddlewareTracing is a client middleware that enables tracing feature using standards of OpenTelemetry.
func MiddlewareTracing(c *Client, r *http.Request) (response *Response, err error) {
	tr := otel.GetTracerProvider().Tracer(
		"github.com/gogf/gf/net/ghttp.Client",
		trace.WithInstrumentationVersion(fmt.Sprintf(`%s`, gf.VERSION)),
	)
	ctx, span := tr.Start(r.Context(), r.URL.String())
	defer span.End()
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
	var resBodyContent string
	if response.ContentLength <= maxContentLogSize {
		reqBodyContentBytes, _ := ioutil.ReadAll(response.Body)
		resBodyContent = string(reqBodyContentBytes)
		response.Body = utils.NewReadCloser(reqBodyContentBytes, false)
	} else {
		resBodyContent = fmt.Sprintf("[Response Body Too Large For Logging, Max: %d bytes]", maxContentLogSize)
	}
	if response != nil {
		span.AddEvent("http.response", trace.WithAttributes(
			label.Any(`http.response.headers`, headerToMap(response.Header)),
			label.String(`http.response.body`, resBodyContent),
		))
	}
	return
}

// headerToMap coverts request headers to map.
func headerToMap(header http.Header) map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range header {
		if len(v) > 1 {
			m[k] = v
		} else {
			m[k] = v[0]
		}
	}
	return m
}
