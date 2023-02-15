// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package response

import (
	"bufio"
	"net"
	"net/http"
)

// Writer wraps http.ResponseWriter for extra features.
type Writer struct {
	http.ResponseWriter      // The underlying ResponseWriter.
	hijacked            bool // Mark this request is hijacked or not.
	wroteHeader         bool // Is header wrote or not, avoiding error: superfluous/multiple response.WriteHeader call.
}

// NewWriter creates and returns a new Writer.
func NewWriter(writer http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: writer,
	}
}

// WriteHeader implements the interface of http.ResponseWriter.WriteHeader.
func (w *Writer) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.wroteHeader = true
}

// Hijack implements the interface function of http.Hijacker.Hijack.
func (w *Writer) Hijack() (conn net.Conn, writer *bufio.ReadWriter, err error) {
	conn, writer, err = w.ResponseWriter.(http.Hijacker).Hijack()
	w.hijacked = true
	return
}

// IsHeaderWrote returns if the header status is written.
func (w *Writer) IsHeaderWrote() bool {
	return w.wroteHeader
}

// IsHijacked returns if the connection is hijacked.
func (w *Writer) IsHijacked() bool {
	return w.hijacked
}

// Flush sends any buffered data to the client.
func (w *Writer) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
		w.wroteHeader = true
	}
}
