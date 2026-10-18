// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests that the client tracing middleware does not read or replace the request body.

package gclient

import (
	"context"
	"io"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/gogf/gf/v2/test/gtest"
)

// readCountReadCloser counts Read calls so a test can see whether the body was consumed.
type readCountReadCloser struct {
	// reads is the number of times Read has been called.
	reads int
}

// Read records one read and returns EOF without producing bytes.
func (r *readCountReadCloser) Read(_ []byte) (int, error) {
	r.reads++
	return 0, io.EOF
}

// Close releases the reader. Closing is not a body read.
func (r *readCountReadCloser) Close() error {
	return nil
}

// Test_ClientTracerTracing_DoesNotReadRequestBody checks that the tracing middleware
// leaves the original request body unread. A later io.ReadAll in the tracer fails this test.
func Test_ClientTracerTracing_DoesNotReadRequestBody(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		originalProvider := otel.GetTracerProvider()
		otel.SetTracerProvider(sdktrace.NewTracerProvider())
		defer otel.SetTracerProvider(originalProvider)

		var (
			client = New()
			body   = &readCountReadCloser{}
		)
		req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1/upload", body)
		t.AssertNil(err)
		// handlerIndex starts at the tracing middleware, so its Next call stops at the stub.
		// The stub must not read req.Body.
		ctx := context.WithValue(req.Context(), clientMiddlewareKey, &clientMiddleware{
			client: client,
			handlers: []HandlerFunc{
				internalMiddlewareObservability,
				func(c *Client, r *http.Request) (*Response, error) {
					return nil, nil
				},
			},
			handlerIndex: 0,
		})
		req = req.WithContext(ctx)

		_, err = internalMiddlewareObservability(client, req)
		t.AssertNil(err)
		t.Assert(body.reads, 0)
		t.Assert(req.Body == body, true)
	})
}
