// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/errors/gerror"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gspath"
	"github.com/gogf/gf/os/gtime"
)

// ServeHTTP is the default handler for http request.
// It should not create new goroutine handling the request as
// it's called by am already created new goroutine from http.Server.
//
// This function also make serve implementing the interface of http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Max body size limit.
	if s.config.ClientMaxBodySize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, s.config.ClientMaxBodySize)
	}

	// Rewrite feature checks.
	if len(s.config.Rewrites) > 0 {
		if rewrite, ok := s.config.Rewrites[r.URL.Path]; ok {
			r.URL.Path = rewrite
		}
	}

	// Remove char '/' in the tail of URI.
	if r.URL.Path != "/" {
		for len(r.URL.Path) > 0 && r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
		}
	}

	// Default URI value if it's empty.
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	// Create a new request object.
	request := newRequest(s, r, w)

	defer func() {
		request.LeaveTime = gtime.TimestampMilli()
		// error log handling.
		if request.error != nil {
			s.handleErrorLog(request.error, request)
		} else {
			if exception := recover(); exception != nil {
				request.Response.WriteStatus(http.StatusInternalServerError)
				s.handleErrorLog(gerror.Newf("%v", exception), request)
			}
		}
		// access log handling.
		s.handleAccessLog(request)
		// Close the session, which automatically update the TTL
		// of the session if it exists.
		request.Session.Close()

		// Close the request and response body
		// to release the file descriptor in time.
		request.Request.Body.Close()
		if request.Request.Response != nil {
			request.Request.Response.Body.Close()
		}
	}()

	// ============================================================
	// Priority:
	// Static File > Dynamic Service > Static Directory
	// ============================================================

	// Search the static file with most high priority,
	// which also handle the index files feature.
	if s.config.FileServerEnabled {
		request.StaticFile = s.searchStaticFile(r.URL.Path)
		if request.StaticFile != nil {
			request.isFileRequest = true
		}
	}

	// Search the dynamic service handler.
	request.handlers, request.hasHookHandler, request.hasServeHandler = s.getHandlersWithCache(request)

	// Check the service type static or dynamic for current request.
	if request.StaticFile != nil && request.StaticFile.IsDir && request.hasServeHandler {
		request.isFileRequest = false
	}

	// HOOK - BeforeServe
	s.callHookHandler(HookBeforeServe, request)

	// Core serving handling.
	if !request.IsExited() {
		if request.isFileRequest {
			// Static file service.
			s.serveFile(request, request.StaticFile)
		} else {
			if len(request.handlers) > 0 {
				// Dynamic service.
				request.Middleware.Next()
			} else {
				if request.StaticFile != nil && request.StaticFile.IsDir {
					// Serve the directory.
					s.serveFile(request, request.StaticFile)
				} else {
					if len(request.Response.Header()) == 0 &&
						request.Response.Status == 0 &&
						request.Response.BufferLength() == 0 {
						request.Response.WriteHeader(http.StatusNotFound)
					}
				}
			}
		}
	}

	// HOOK - AfterServe
	if !request.IsExited() {
		s.callHookHandler(HookAfterServe, request)
	}

	// HOOK - BeforeOutput
	if !request.IsExited() {
		s.callHookHandler(HookBeforeOutput, request)
	}

	// HTTP status checking.
	if request.Response.Status == 0 {
		if request.StaticFile != nil || request.Middleware.served || request.Response.buffer.Len() > 0 {
			request.Response.WriteHeader(http.StatusOK)
		} else {
			request.Response.WriteHeader(http.StatusNotFound)
		}
	}
	// HTTP status handler.
	if request.Response.Status != http.StatusOK {
		statusFuncArray := s.getStatusHandler(request.Response.Status, request)
		for _, f := range statusFuncArray {
			// Call custom status handler.
			niceCallFunc(func() {
				f(request)
			})
			if request.IsExited() {
				break
			}
		}
	}

	// Automatically set the session id to cookie
	// if it creates a new session id in this request.
	if s.config.SessionCookieOutput &&
		request.Session.IsDirty() &&
		request.Session.Id() != request.GetSessionId() {
		request.Cookie.SetSessionId(request.Session.Id())
	}
	// Output the cookie content to client.
	request.Cookie.Flush()
	// Output the buffer content to client.
	request.Response.Flush()
	// HOOK - AfterOutput
	if !request.IsExited() {
		s.callHookHandler(HookAfterOutput, request)
	}
}

// searchStaticFile searches the file with given URI.
// It returns a file struct specifying the file information.
func (s *Server) searchStaticFile(uri string) *StaticFile {
	var file *gres.File
	var path string
	var dir bool
	// Firstly search the StaticPaths mapping.
	if len(s.config.StaticPaths) > 0 {
		for _, item := range s.config.StaticPaths {
			if len(uri) >= len(item.prefix) && strings.EqualFold(item.prefix, uri[0:len(item.prefix)]) {
				// To avoid case like: /static/style -> /static/style.css
				if len(uri) > len(item.prefix) && uri[len(item.prefix)] != '/' {
					continue
				}
				file = gres.GetWithIndex(item.path+uri[len(item.prefix):], s.config.IndexFiles)
				if file != nil {
					return &StaticFile{
						File:  file,
						IsDir: file.FileInfo().IsDir(),
					}
				}
				path, dir = gspath.Search(item.path, uri[len(item.prefix):], s.config.IndexFiles...)
				if path != "" {
					return &StaticFile{
						Path:  path,
						IsDir: dir,
					}
				}

			}
		}
	}
	// Secondly search the root and searching paths.
	if len(s.config.SearchPaths) > 0 {
		for _, p := range s.config.SearchPaths {
			file = gres.GetWithIndex(p+uri, s.config.IndexFiles)
			if file != nil {
				return &StaticFile{
					File:  file,
					IsDir: file.FileInfo().IsDir(),
				}
			}
			if path, dir = gspath.Search(p, uri, s.config.IndexFiles...); path != "" {
				return &StaticFile{
					Path:  path,
					IsDir: dir,
				}
			}
		}
	}
	// Lastly search the resource manager.
	if len(s.config.StaticPaths) == 0 && len(s.config.SearchPaths) == 0 {
		if file = gres.GetWithIndex(uri, s.config.IndexFiles); file != nil {
			return &StaticFile{
				File:  file,
				IsDir: file.FileInfo().IsDir(),
			}
		}
	}
	return nil
}

// serveFile serves the static file for client.
// The optional parameter <allowIndex> specifies if allowing directory listing if <f> is directory.
func (s *Server) serveFile(r *Request, f *StaticFile, allowIndex ...bool) {
	// Use resource file from memory.
	if f.File != nil {
		if f.IsDir {
			if s.config.IndexFolder || (len(allowIndex) > 0 && allowIndex[0]) {
				s.listDir(r, f.File)
			} else {
				r.Response.WriteStatus(http.StatusForbidden)
			}
		} else {
			info := f.File.FileInfo()
			r.Response.wroteHeader = true
			http.ServeContent(r.Response.Writer.RawWriter(), r.Request, info.Name(), info.ModTime(), f.File)
		}
		return
	}
	// Use file from dist.
	file, err := os.Open(f.Path)
	if err != nil {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}
	defer file.Close()

	// Clear the response buffer before file serving.
	// It ignores all custom buffer content and uses the file content.
	r.Response.ClearBuffer()

	info, _ := file.Stat()
	if info.IsDir() {
		if s.config.IndexFolder || (len(allowIndex) > 0 && allowIndex[0]) {
			s.listDir(r, file)
		} else {
			r.Response.WriteStatus(http.StatusForbidden)
		}
	} else {
		r.Response.wroteHeader = true
		http.ServeContent(r.Response.Writer.RawWriter(), r.Request, info.Name(), info.ModTime(), file)
	}
}

// listDir lists the sub files of specified directory as HTML content to client.
func (s *Server) listDir(r *Request, f http.File) {
	files, err := f.Readdir(-1)
	if err != nil {
		r.Response.WriteStatus(http.StatusInternalServerError, "Error reading directory")
		return
	}
	// The folder type has the most priority than file.
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].Name() < files[j].Name()
	})
	if r.Response.Header().Get("Content-Type") == "" {
		r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	r.Response.Write(`<html>`)
	r.Response.Write(`<head>`)
	r.Response.Write(`<style>`)
	r.Response.Write(`body {font-family:Consolas, Monaco, "Andale Mono", "Ubuntu Mono", monospace;}`)
	r.Response.Write(`</style>`)
	r.Response.Write(`</head>`)
	r.Response.Write(`<body>`)
	r.Response.Writef(`<h1>Index of %s</h1>`, r.URL.Path)
	r.Response.Writef(`<hr />`)
	r.Response.Write(`<table>`)
	if r.URL.Path != "/" {
		r.Response.Write(`<tr>`)
		r.Response.Writef(`<td><a href="%s">..</a></td>`, gfile.Dir(r.URL.Path))
		r.Response.Write(`</tr>`)
	}
	name := ""
	size := ""
	prefix := gstr.TrimRight(r.URL.Path, "/")
	for _, file := range files {
		name = file.Name()
		size = gfile.FormatSize(file.Size())
		if file.IsDir() {
			name += "/"
			size = "-"
		}
		r.Response.Write(`<tr>`)
		r.Response.Writef(`<td><a href="%s/%s">%s</a></td>`, prefix, name, ghtml.SpecialChars(name))
		r.Response.Writef(`<td style="width:300px;text-align:center;">%s</td>`, gtime.New(file.ModTime()).ISO8601())
		r.Response.Writef(`<td style="width:80px;text-align:right;">%s</td>`, size)
		r.Response.Write(`</tr>`)
	}
	r.Response.Write(`</table>`)
	r.Response.Write(`</body>`)
	r.Response.Write(`</html>`)
}
