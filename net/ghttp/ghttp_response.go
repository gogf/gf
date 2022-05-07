// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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
	"net/url"

	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
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
	var serveFile *staticFile
	if file := gres.Get(path); file != nil {
		serveFile = &staticFile{
			File:  file,
			IsDir: file.FileInfo().IsDir(),
		}
	} else {
		path, _ = gfile.Search(path)
		if path == "" {
			r.WriteStatus(http.StatusNotFound)
			return
		}
		serveFile = &staticFile{Path: path}
	}
	r.Server.serveFile(r.Request, serveFile, allowIndex...)
}

// ServeFileDownload serves file downloading to the response.
func (r *Response) ServeFileDownload(path string, name ...string) {
	var (
		serveFile    *staticFile
		downloadName = ""
	)

	if len(name) > 0 {
		downloadName = name[0]
	}
	if file := gres.Get(path); file != nil {
		serveFile = &staticFile{
			File:  file,
			IsDir: file.FileInfo().IsDir(),
		}
		if downloadName == "" {
			downloadName = gfile.Basename(file.Name())
		}
	} else {
		path, _ = gfile.Search(path)
		if path == "" {
			r.WriteStatus(http.StatusNotFound)
			return
		}
		serveFile = &staticFile{Path: path}
		if downloadName == "" {
			downloadName = gfile.Basename(path)
		}
	}
	r.Header().Set("Content-Type", "application/force-download")
	r.Header().Set("Accept-Ranges", "bytes")
	r.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename=%s`, url.QueryEscape(downloadName)))
	r.Server.serveFile(r.Request, serveFile)
}

// RedirectTo redirects the client to another location.
// The optional parameter `code` specifies the http status code for redirecting,
// which commonly can be 301 or 302. It's 302 in default.
func (r *Response) RedirectTo(location string, code ...int) {
	r.Header().Set("Location", location)
	if len(code) > 0 {
		r.WriteHeader(code[0])
	} else {
		r.WriteHeader(http.StatusFound)
	}
	r.Request.Exit()
}

// RedirectBack redirects the client back to referer.
// The optional parameter `code` specifies the http status code for redirecting,
// which commonly can be 301 or 302. It's 302 in default.
func (r *Response) RedirectBack(code ...int) {
	r.RedirectTo(r.Request.GetReferer(), code...)
}

// Buffer returns the buffered content as []byte.
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

// SetBuffer overwrites the buffer with `data`.
func (r *Response) SetBuffer(data []byte) {
	r.buffer.Reset()
	r.buffer.Write(data)
}

// ClearBuffer clears the response buffer.
func (r *Response) ClearBuffer() {
	r.buffer.Reset()
}

// Flush outputs the buffer content to the client and clears the buffer.
func (r *Response) Flush() {
	r.Header().Set(responseTraceIDHeader, gtrace.GetTraceID(r.Request.Context()))
	if r.Server.config.ServerAgent != "" {
		r.Header().Set("Server", r.Server.config.ServerAgent)
	}
	r.Writer.Flush()
}
