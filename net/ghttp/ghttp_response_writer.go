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

// ResponseWriter is the custom writer for http response.
type ResponseWriter struct {
	Status      int                 // HTTP status.
	writer      http.ResponseWriter // The underlying ResponseWriter.
	buffer      *bytes.Buffer       // The output buffer.
	hijacked    bool                // Mark this request is hijacked or not.
	wroteHeader bool                // Is header wrote or not, avoiding error: superfluous/multiple response.WriteHeader call.
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
	w.hijacked = true
	return w.writer.(http.Hijacker).Hijack()
}

// OutputBuffer outputs the buffer to client and clears the buffer.
func (w *ResponseWriter) Flush() {
	if w.hijacked {
		return
	}
	if w.Status != 0 && !w.wroteHeader {
		w.wroteHeader = true
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
