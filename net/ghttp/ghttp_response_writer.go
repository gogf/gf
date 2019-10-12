// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

// Custom ResponseWriter, which is used for controlling the output buffer.
type ResponseWriter struct {
	Status      int                 // HTTP status.
	writer      http.ResponseWriter // The underlying ResponseWriter.
	buffer      *bytes.Buffer       // The output buffer.
	wroteHeader bool                // Is header wrote, avoiding error: superfluous/multiple response.WriteHeader call.
}

// RawWriter returns the underlying ResponseWriter.
func (w *ResponseWriter) RawWriter() http.ResponseWriter {
	return w.writer
}

// Header implements the interface function of http.ResponseWriter.Header.
func (w *ResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write implements the interface function of http.ResponseWriter.Write.
func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)
	return len(data), nil
}

// WriteHeader implements the interface of http.ResponseWriter.WriteHeader.
func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
}

// Hijack implements the interface function of http.Hijacker.Hijack.
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.writer.(http.Hijacker).Hijack()
}

// OutputBuffer outputs the buffer to client.
func (w *ResponseWriter) OutputBuffer() {
	if w.Status != 0 && !w.wroteHeader {
		w.writer.WriteHeader(w.Status)
	}
	// Default status text output.
	if w.Status != http.StatusOK && w.buffer.Len() == 0 {
		w.buffer.WriteString(http.StatusText(w.Status))
	}
	if w.buffer.Len() > 0 {
		w.writer.Write(w.buffer.Bytes())
		w.buffer.Reset()
	}
}
