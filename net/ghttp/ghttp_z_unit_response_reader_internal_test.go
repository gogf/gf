// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogf/gf/v2/net/ghttp/internal/response"
	"github.com/gogf/gf/v2/test/gtest"
)

type countingReader struct {
	total int
	read  int
}

func (r *countingReader) Read(p []byte) (int, error) {
	if r.read >= r.total {
		return 0, io.EOF
	}
	n := len(p)
	if remaining := r.total - r.read; n > remaining {
		n = remaining
	}
	for i := 0; i < n; i++ {
		p[i] = 'x'
	}
	r.read += n
	return n, nil
}

type failingResponseWriter struct {
	header  http.Header
	limit   int
	written int
}

func (w *failingResponseWriter) Header() http.Header {
	return w.header
}

func (w *failingResponseWriter) WriteHeader(int) {}

func (w *failingResponseWriter) Write(p []byte) (int, error) {
	if w.written >= w.limit {
		return 0, errors.New("client disconnected")
	}
	w.written += len(p)
	return len(p), nil
}

func Test_Response_ServeReaderWriteFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			reader = &countingReader{total: 1 << 20}
			writer = &failingResponseWriter{header: make(http.Header), limit: 32 * 1024}
			resp   = &Response{
				BufferWriter: response.NewBufferWriter(writer),
				Request:      &Request{Request: httptest.NewRequest(http.MethodGet, "/", nil)},
			}
		)
		resp.ServeReader(reader, int64(reader.total))
		// io.Copy must stop once writing fails instead of draining the whole reader.
		t.AssertLT(reader.read, reader.total)
		t.Assert(writer.Header().Get(HeaderContentLength), "1048576")
	})
}
