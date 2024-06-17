// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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

	"github.com/gogf/gf/v2/encoding/ghtml"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gspath"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
)

// ServeHTTP is the default handler for http request.
// It should not create new goroutine handling the request as
// it's called by am already created new goroutine from http.Server.
//
// This function also makes serve implementing the interface of http.Handler.
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

	var (
		request   = newRequest(s, r, w)    // Create a new request object.
		sessionId = request.GetSessionId() // Get sessionId before user handler
	)
	defer s.handleAfterRequestDone(request)

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
	request.handlers,
		request.serveHandler,
		request.hasHookHandler,
		request.hasServeHandler = s.getHandlersWithCache(request)

	// Check the service type static or dynamic for current request.
	if request.StaticFile != nil && request.StaticFile.IsDir && request.hasServeHandler {
		request.isFileRequest = false
	}

	// Metrics.
	s.handleMetricsBeforeRequest(request)

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

	// Response handling.
	s.handleResponse(request, sessionId)

	// HOOK - AfterOutput
	if !request.IsExited() {
		s.callHookHandler(HookAfterOutput, request)
	}
}

func (s *Server) handleResponse(request *Request, sessionId string) {
	// HTTP status checking.
	if request.Response.Status == 0 {
		if request.StaticFile != nil || request.Middleware.served || request.Response.BufferLength() > 0 {
			request.Response.WriteHeader(http.StatusOK)
		} else if err := request.GetError(); err != nil {
			if request.Response.BufferLength() == 0 {
				request.Response.Write(err.Error())
			}
			request.Response.WriteHeader(http.StatusInternalServerError)
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
	// if it creates a new session id in this request
	// and SessionCookieOutput is enabled.
	if s.config.SessionCookieOutput && request.Session.IsDirty() {
		// Can change by r.Session.SetId("") before init session
		// Can change by r.Cookie.SetSessionId("")
		sidFromSession, sidFromRequest := request.Session.MustId(), request.GetSessionId()
		if sidFromSession != sidFromRequest {
			if sidFromSession != sessionId {
				request.Cookie.SetSessionId(sidFromSession)
			} else {
				request.Cookie.SetSessionId(sidFromRequest)
			}
		}
	}
	// Output the cookie content to the client.
	request.Cookie.Flush()
	// Output the buffer content to the client.
	request.Response.Flush()
}

func (s *Server) handleAfterRequestDone(request *Request) {
	request.LeaveTime = gtime.Now()
	// error log handling.
	if request.error != nil {
		s.handleErrorLog(request.error, request)
	} else {
		if exception := recover(); exception != nil {
			request.Response.WriteStatus(http.StatusInternalServerError)
			if v, ok := exception.(error); ok {
				if code := gerror.Code(v); code != gcode.CodeNil {
					s.handleErrorLog(v, request)
				} else {
					s.handleErrorLog(
						gerror.WrapCodeSkip(gcode.CodeInternalPanic, 1, v, ""),
						request,
					)
				}
			} else {
				s.handleErrorLog(
					gerror.NewCodeSkipf(gcode.CodeInternalPanic, 1, "%+v", exception),
					request,
				)
			}
		}
	}
	// access log handling.
	s.handleAccessLog(request)
	// Close the session, which automatically update the TTL
	// of the session if it exists.
	if err := request.Session.Close(); err != nil {
		intlog.Errorf(request.Context(), `%+v`, err)
	}

	// Close the request and response body
	// to release the file descriptor in time.
	err := request.Request.Body.Close()
	if err != nil {
		intlog.Errorf(request.Context(), `%+v`, err)
	}
	if request.Request.Response != nil {
		err = request.Request.Response.Body.Close()
		if err != nil {
			intlog.Errorf(request.Context(), `%+v`, err)
		}
	}

	// Metrics.
	s.handleMetricsAfterRequestDone(request)
}

// searchStaticFile searches the file with given URI.
// It returns a file struct specifying the file information.
func (s *Server) searchStaticFile(uri string) *staticFile {
	var (
		file *gres.File
		path string
		dir  bool
	)
	// Firstly search the StaticPaths mapping.
	if len(s.config.StaticPaths) > 0 {
		for _, item := range s.config.StaticPaths {
			if len(uri) >= len(item.Prefix) && strings.EqualFold(item.Prefix, uri[0:len(item.Prefix)]) {
				// To avoid case like: /static/style -> /static/style.css
				if len(uri) > len(item.Prefix) && uri[len(item.Prefix)] != '/' {
					continue
				}
				file = gres.GetWithIndex(item.Path+uri[len(item.Prefix):], s.config.IndexFiles)
				if file != nil {
					return &staticFile{
						File:  file,
						IsDir: file.FileInfo().IsDir(),
					}
				}
				path, dir = gspath.Search(item.Path, uri[len(item.Prefix):], s.config.IndexFiles...)
				if path != "" {
					return &staticFile{
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
				return &staticFile{
					File:  file,
					IsDir: file.FileInfo().IsDir(),
				}
			}
			if path, dir = gspath.Search(p, uri, s.config.IndexFiles...); path != "" {
				return &staticFile{
					Path:  path,
					IsDir: dir,
				}
			}
		}
	}
	// Lastly search the resource manager.
	if len(s.config.StaticPaths) == 0 && len(s.config.SearchPaths) == 0 {
		if file = gres.GetWithIndex(uri, s.config.IndexFiles); file != nil {
			return &staticFile{
				File:  file,
				IsDir: file.FileInfo().IsDir(),
			}
		}
	}
	return nil
}

// serveFile serves the static file for the client.
// The optional parameter `allowIndex` specifies if allowing directory listing if `f` is a directory.
func (s *Server) serveFile(r *Request, f *staticFile, allowIndex ...bool) {
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
			r.Response.ServeContent(info.Name(), info.ModTime(), f.File)
		}
		return
	}
	// Use file from dist.
	file, err := os.Open(f.Path)
	if err != nil {
		r.Response.WriteStatus(http.StatusForbidden)
		return
	}
	defer func() {
		_ = file.Close()
	}()

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
		r.Response.ServeContent(info.Name(), info.ModTime(), file)
	}
}

// listDir lists the sub files of specified directory as HTML content to the client.
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
