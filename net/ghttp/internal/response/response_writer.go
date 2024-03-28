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
	http.ResponseWriter       // The underlying ResponseWriter.
	hijacked            bool  // Mark this request is hijacked or not.
	wroteHeader         bool  // Is header wrote or not, avoiding error: superfluous/multiple response.WriteHeader call.
	bytesWritten        int64 // Bytes written to response.
}

// NewWriter creates and returns a new Writer.
func NewWriter(writer http.ResponseWriter) *Writer {
	return &Writer{
		ResponseWriter: writer,
	}
}

// WriteHeader implements the interface of http.ResponseWriter.WriteHeader.
// Note that the underlying `WriteHeader` can only be called once in a http response.
func (w *Writer) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.ResponseWriter.WriteHeader(status)
	w.wroteHeader = true
}

// BytesWritten returns the length that was written to response.
func (w *Writer) BytesWritten() int64 {
	return w.bytesWritten
}

// Write implements the interface function of http.ResponseWriter.Write.
func (w *Writer) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.bytesWritten += int64(n)
	w.wroteHeader = true
	return n, err
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
