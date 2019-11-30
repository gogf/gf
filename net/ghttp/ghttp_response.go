// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/os/gfile"
)

// Response is the http response manager.
// Note that it implements the http.ResponseWriter interface with buffering feature.
type Response struct {
	*ResponseWriter                 // Underlying ResponseWriter.
	Server          *Server         // Parent server.
	Writer          *ResponseWriter // Alias of ResponseWriter.
	Request         *Request        // According request.
}

// newResponse creates and returns a new Response object.
func newResponse(s *Server, w http.ResponseWriter) *Response {
	r := &Response{
		Server: s,
		ResponseWriter: &ResponseWriter{
			writer: w,
			buffer: bytes.NewBuffer(nil),
		},
	}
	r.Writer = r.ResponseWriter
	return r
}

// ServeFile serves the file to the response.
func (r *Response) ServeFile(path string, allowIndex ...bool) {
	serveFile := (*staticServeFile)(nil)
	if file := gres.Get(path); file != nil {
		serveFile = &staticServeFile{
			file: file,
			dir:  file.FileInfo().IsDir(),
		}
	} else {
		path = gfile.RealPath(path)
		if path == "" {
			r.WriteStatus(http.StatusNotFound)
			return
		}
		serveFile = &staticServeFile{path: path}
	}
	r.Server.serveFile(r.Request, serveFile, allowIndex...)
}

// ServeFileDownload serves file downloading to the response.
func (r *Response) ServeFileDownload(path string, name ...string) {
	serveFile := (*staticServeFile)(nil)
	downloadName := ""
	if len(name) > 0 {
		downloadName = name[0]
	}
	if file := gres.Get(path); file != nil {
		serveFile = &staticServeFile{
			file: file,
			dir:  file.FileInfo().IsDir(),
		}
		if downloadName == "" {
			downloadName = gfile.Basename(file.Name())
		}
	} else {
		path = gfile.RealPath(path)
		if path == "" {
			r.WriteStatus(http.StatusNotFound)
			return
		}
		serveFile = &staticServeFile{path: path}
		if downloadName == "" {
			downloadName = gfile.Basename(path)
		}
	}
	r.Header().Set("Content-Type", "application/force-download")
	r.Header().Set("Accept-Ranges", "bytes")
	r.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, downloadName))
	r.Server.serveFile(r.Request, serveFile)
}

// RedirectTo redirects client to another location using http status 302.
func (r *Response) RedirectTo(location string) {
	r.Header().Set("Location", location)
	r.WriteHeader(http.StatusFound)
	r.Request.Exit()
}

// RedirectBack redirects client back to referer using http status 302.
func (r *Response) RedirectBack() {
	r.RedirectTo(r.Request.GetReferer())
}

// BufferString returns the buffered content as []byte.
func (r *Response) Buffer() []byte {
	return r.buffer.Bytes()
}

// BufferString returns the buffered content as string.
func (r *Response) BufferString() string {
	return r.buffer.String()
}

// BufferLength returns the length of the buffered content.
func (r *Response) BufferLength() int {
	return r.buffer.Len()
}

// SetBuffer overwrites the buffer with <data>.
func (r *Response) SetBuffer(data []byte) {
	r.buffer.Reset()
	r.buffer.Write(data)
}

// ClearBuffer clears the response buffer.
func (r *Response) ClearBuffer() {
	r.buffer.Reset()
}

// Output outputs the buffer content to the client.
func (r *Response) Output() {
	if r.Server.config.ServerAgent != "" {
		r.Header().Set("Server", r.Server.config.ServerAgent)
	}
	r.Writer.OutputBuffer()
}
