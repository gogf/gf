// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"io"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/ghttp/internal/response"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

// Response is the http response manager.
// Note that it implements the http.ResponseWriter interface with buffering feature.
type Response struct {
	*response.BufferWriter          // Underlying ResponseWriter.
	Server                 *Server  // Parent server.
	Request                *Request // According request.
}

// newResponse creates and returns a new Response object.
func newResponse(s *Server, w http.ResponseWriter) *Response {
	r := &Response{
		Server:       s,
		BufferWriter: response.NewBufferWriter(w),
	}
	return r
}

// ServeFile serves the file to the response.
func (r *Response) ServeFile(path string, allowIndex ...bool) {
	var (
		serveFile *staticFile
	)
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
	r.Header().Set(HeaderContentDisposition, formatContentDisposition(downloadName))
	r.Header().Set(HeaderAccessControlExposeHeaders, HeaderContentDisposition)
	r.Server.serveFile(r.Request, serveFile)
}

// ServeReader streams the content of reader to the client without buffering it
// in memory. The size specifies the number of bytes to send, and a negative
// size means the content length is unknown. Range requests are not supported
// and no Accept-Ranges header is set. If reader implements io.Closer, it is
// closed after the content is sent.
//
// Note that the content is written to the client directly, so it should not be
// mixed with Response.Write in the same handler.
func (r *Response) ServeReader(reader io.Reader, size int64) {
	r.serveReader(reader, size)
}

// ServeReaderDownload streams the content of reader as a download. The
// contentType defaults to application/octet-stream when omitted. Its size,
// Range and io.Closer behavior is the same as ServeReader.
func (r *Response) ServeReaderDownload(reader io.Reader, size int64, name string, contentType ...string) {
	r.Header().Set(HeaderContentDisposition, formatContentDisposition(name))
	r.Header().Set(HeaderAccessControlExposeHeaders, HeaderContentDisposition)
	downloadContentType := "application/octet-stream"
	if len(contentType) > 0 && contentType[0] != "" {
		downloadContentType = contentType[0]
	}
	r.Header().Set(HeaderContentType, downloadContentType)
	r.serveReader(reader, size)
}

// serveReader is the shared implementation of ServeReader and ServeReaderDownload.
func (r *Response) serveReader(reader io.Reader, size int64) {
	if size >= 0 {
		r.Header().Set(HeaderContentLength, strconv.FormatInt(size, 10))
	}
	if closer, ok := reader.(io.Closer); ok {
		defer func() {
			if err := closer.Close(); err != nil {
				intlog.Errorf(r.Request.Context(), `close response reader failed: %+v`, err)
			}
		}()
	}
	writer := r.RawWriter()
	// HEAD requests receive the headers only and must not consume the reader.
	if r.Request.Method == http.MethodHead {
		writer.WriteHeader(http.StatusOK)
		return
	}
	if _, err := io.Copy(writer, reader); err != nil {
		intlog.Errorf(r.Request.Context(), `stream response content failed: %+v`, err)
	}
}

// formatContentDisposition returns an RFC 6266 attachment value with a sanitized
// filename parameter.
func formatContentDisposition(name string) string {
	// Normalize backslashes first so paths are split the same way on all platforms.
	name = path.Base(strings.ReplaceAll(strings.TrimSpace(name), `\`, "/"))
	name = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, name)
	name = strings.Trim(strings.TrimSpace(name), "/")
	if name == "" || name == "." || name == ".." {
		name = "download"
	}
	return mime.FormatMediaType("attachment", map[string]string{"filename": name})
}

// RedirectTo redirects the client to another location.
// The optional parameter `code` specifies the http status code for redirecting,
// which commonly can be 301 or 302. It's 302 in default.
func (r *Response) RedirectTo(location string, code ...int) {
	r.Header().Set(HeaderLocation, location)
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

// ServeContent replies to the request using the content in the
// provided ReadSeeker. The main benefit of ServeContent over io.Copy
// is that it handles Range requests properly, sets the MIME type, and
// handles If-Match, If-Unmodified-Since, If-None-Match, If-Modified-Since,
// and If-Range requests.
//
// See http.ServeContent
func (r *Response) ServeContent(name string, modTime time.Time, content io.ReadSeeker) {
	http.ServeContent(r.RawWriter(), r.Request.Request, name, modTime, content)
}

// Flush outputs the buffer content to the client and clears the buffer.
func (r *Response) Flush() {
	r.Header().Set(HeaderTraceId, gtrace.GetTraceID(r.Request.Context()))
	if r.Server.config.ServerAgent != "" {
		r.Header().Set(HeaderServer, r.Server.config.ServerAgent)
	}
	r.BufferWriter.Flush()
}
