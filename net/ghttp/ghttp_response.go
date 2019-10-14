// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/util/gconv"
)

// Response is the writer for response buffer.
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

// Write writes <content> to the response buffer.
func (r *Response) Write(content ...interface{}) {
	if len(content) == 0 {
		return
	}
	if r.Status == 0 {
		r.Status = http.StatusOK
	}
	for _, v := range content {
		switch value := v.(type) {
		case []byte:
			r.buffer.Write(value)
		case string:
			r.buffer.WriteString(value)
		default:
			r.buffer.WriteString(gconv.String(v))
		}
	}
}

// WriteOver overwrites the response buffer with <content>.
func (r *Response) WriteOver(content ...interface{}) {
	r.ClearBuffer()
	r.Write(content...)
}

// Writef writes the response with fmt.Sprintf.
func (r *Response) Writef(format string, params ...interface{}) {
	r.Write(fmt.Sprintf(format, params...))
}

// Writef writes the response with <content> and new line.
func (r *Response) Writeln(content ...interface{}) {
	if len(content) == 0 {
		r.Write("\n")
		return
	}
	r.Write(append(content, "\n")...)
}

// Writefln writes the response with fmt.Sprintf and new line.
func (r *Response) Writefln(format string, params ...interface{}) {
	r.Writeln(fmt.Sprintf(format, params...))
}

// WriteJson writes <content> to the response with JSON format.
func (r *Response) WriteJson(content interface{}) error {
	if b, err := json.Marshal(content); err != nil {
		return err
	} else {
		r.Header().Set("Content-Type", "application/json")
		r.Write(b)
	}
	return nil
}

// WriteJson writes <content> to the response with JSONP format.
// Note that there should be a "callback" parameter in the request for JSONP format.
func (r *Response) WriteJsonP(content interface{}) error {
	if b, err := json.Marshal(content); err != nil {
		return err
	} else {
		//r.Header().Set("Content-Type", "application/json")
		if callback := r.Request.GetString("callback"); callback != "" {
			buffer := []byte(callback)
			buffer = append(buffer, byte('('))
			buffer = append(buffer, b...)
			buffer = append(buffer, byte(')'))
			r.Write(buffer)
		} else {
			r.Write(b)
		}
	}
	return nil
}

// WriteJson writes <content> to the response with XML format.
func (r *Response) WriteXml(content interface{}, rootTag ...string) error {
	if b, err := gparser.VarToXml(content, rootTag...); err != nil {
		return err
	} else {
		r.Header().Set("Content-Type", "application/xml")
		r.Write(b)
	}
	return nil
}

// WriteStatus writes HTTP <status> and <content> to the response.
func (r *Response) WriteStatus(status int, content ...interface{}) {
	r.WriteHeader(status)
	if len(content) > 0 {
		r.Write(content...)
	} else {
		r.Write(http.StatusText(status))
	}
	if r.Header().Get("Content-Type") == "" {
		r.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//r.Header().Set("X-Content-Type-Options", "nosniff")
	}
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

// ServeFileDownload serves file as file downloading to the response.
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

// RedirectTo redirects client to another location.
func (r *Response) RedirectTo(location string) {
	r.Header().Set("Location", location)
	r.WriteHeader(http.StatusFound)
	r.Request.Exit()
}

// RedirectBack redirects client back to referer.
func (r *Response) RedirectBack() {
	r.RedirectTo(r.Request.GetReferer())
}

// BufferString returns the buffer content as []byte.
func (r *Response) Buffer() []byte {
	return r.buffer.Bytes()
}

// BufferString returns the buffer content as string.
func (r *Response) BufferString() string {
	return r.buffer.String()
}

// BufferLength returns the length of the buffer content.
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
	r.Header().Set("Server", r.Server.config.ServerAgent)
	//r.handleGzip()
	r.Writer.OutputBuffer()
}
