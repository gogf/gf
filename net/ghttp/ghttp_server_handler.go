// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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

// 服务静态文件信息
type staticServeFile struct {
	file *gres.File // 资源文件
	path string     // 文件路径
	dir  bool       // 是否目录
}

// 默认HTTP Server处理入口，http包底层默认使用了gorutine异步处理请求，所以这里不再异步执行
func (s *Server) defaultHttpHandle(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r)
}

// 执行处理HTTP请求，
// 首先，查找是否有对应域名的处理接口配置；
// 其次，如果没有对应的自定义处理接口配置，那么走默认的域名处理接口配置；
// 最后，如果以上都没有找到处理接口，那么进行文件处理；
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	// 重写规则判断
	if len(s.config.Rewrites) > 0 {
		if rewrite, ok := s.config.Rewrites[r.URL.Path]; ok {
			r.URL.Path = rewrite
		}
	}

	// URI默认值
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	// 去掉末尾的"/"号
	if r.URL.Path != "/" {
		for r.URL.Path[len(r.URL.Path)-1] == '/' {
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
		}
	}

	// 创建请求处理对象
	request := newRequest(s, r, w)

	defer func() {
		// 设置请求完成时间
		request.LeaveTime = gtime.Microsecond()
		// error log
		if request.error != nil {
			s.handleErrorLog(request.error, request)
		} else {
			if exception := recover(); exception != nil {
				request.Response.WriteStatus(http.StatusInternalServerError)
				s.handleErrorLog(gerror.Newf("%v", exception), request)
			}
		}
		// access log
		s.handleAccessLog(request)
	}()

	// ============================================================
	// 优先级控制:
	// 静态文件 > 动态服务 > 静态目录
	// ============================================================

	serveFile := (*staticServeFile)(nil)
	// 优先执行静态文件检索(检测是否存在对应的静态文件，包括index files处理)
	if s.config.FileServerEnabled {
		serveFile = s.searchStaticFile(r.URL.Path)
		if serveFile != nil {
			request.isFileRequest = true
		}
	}

	// 动态服务检索
	request.handlers, request.hasHookHandler, request.hasServeHandler = s.getHandlersWithCache(request)

	// 判断最终对该请求提供的服务方式
	if serveFile != nil && serveFile.dir && request.hasServeHandler {
		request.isFileRequest = false
	}

	// 事件 - BeforeServe
	s.callHookHandler(HOOK_BEFORE_SERVE, request)

	// 执行静态文件服务/回调控制器/执行对象/方法
	if !request.IsExited() {
		if request.isFileRequest {
			// 静态服务
			s.serveFile(request, serveFile)
		} else {
			if len(request.handlers) > 0 {
				// 动态服务
				request.Middleware.Next()
			} else {
				if serveFile != nil && serveFile.dir {
					// 静态目录
					s.serveFile(request, serveFile)
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

	// 事件 - AfterServe
	if !request.IsExited() {
		s.callHookHandler(HOOK_AFTER_SERVE, request)
	}

	// 事件 - BeforeOutput
	if !request.IsExited() {
		s.callHookHandler(HOOK_BEFORE_OUTPUT, request)
	}

	// HTTP status checking.
	if request.Response.Status == 0 {
		if request.Middleware.served || request.Response.buffer.Len() > 0 {
			request.Response.WriteHeader(http.StatusOK)
		} else {
			request.Response.WriteHeader(http.StatusNotFound)
		}
	}
	// HTTP status handler.
	if request.Response.Status != http.StatusOK {
		if f := s.getStatusHandler(request.Response.Status, request); f != nil {
			// Call custom status handler.
			niceCallFunc(func() {
				f(request)
			})
		}
	}

	// 设置Session Id到Cookie中
	if request.Session.IsDirty() && request.Session.Id() != request.GetSessionId() {
		request.Cookie.SetSessionId(request.Session.Id())
	}
	// 输出Cookie
	request.Cookie.Output()
	// 输出缓冲区
	request.Response.Output()
	// 事件 - AfterOutput
	if !request.IsExited() {
		s.callHookHandler(HOOK_AFTER_OUTPUT, request)
	}
	// 关闭当前Session，并更新会话超时时间
	request.Session.Close()
}

// 查找静态文件的绝对路径
func (s *Server) searchStaticFile(uri string) *staticServeFile {
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
				// Firstly searching resource manager.
				file = gres.GetWithIndex(item.path+uri[len(item.prefix):], s.config.IndexFiles)
				if file != nil {
					return &staticServeFile{
						file: file,
						dir:  file.FileInfo().IsDir(),
					}
				}
				// Secondly searching the file system.
				path, dir = gspath.Search(item.path, uri[len(item.prefix):], s.config.IndexFiles...)
				if path != "" {
					return &staticServeFile{
						path: path,
						dir:  dir,
					}
				}

			}
		}
	}
	// 其次查找root和search path
	if len(s.config.SearchPaths) > 0 {
		for _, p := range s.config.SearchPaths {
			// 优先检索资源管理器
			file = gres.GetWithIndex(p+uri, s.config.IndexFiles)
			if file != nil {
				return &staticServeFile{
					file: file,
					dir:  file.FileInfo().IsDir(),
				}
			}
			// 其次检索文件系统
			if path, dir = gspath.Search(p, uri, s.config.IndexFiles...); path != "" {
				return &staticServeFile{
					path: path,
					dir:  dir,
				}
			}
		}
	}
	// 最后通过资源对象+URI进行文件检索
	if len(s.config.StaticPaths) == 0 && len(s.config.SearchPaths) == 0 {
		if file = gres.GetWithIndex(uri, s.config.IndexFiles); file != nil {
			return &staticServeFile{
				file: file,
				dir:  file.FileInfo().IsDir(),
			}
		}
	}
	return nil
}

// http server静态文件处理，path可以为相对路径也可以为绝对路径
func (s *Server) serveFile(r *Request, f *staticServeFile, allowIndex ...bool) {
	// 使用资源文件
	if f.file != nil {
		if f.dir {
			if s.config.IndexFolder || (len(allowIndex) > 0 && allowIndex[0]) {
				s.listDir(r, f.file)
			} else {
				r.Response.WriteStatus(http.StatusForbidden)
			}
		} else {
			info := f.file.FileInfo()
			r.Response.wroteHeader = true
			http.ServeContent(r.Response.Writer.RawWriter(), r.Request, info.Name(), info.ModTime(), f.file)
		}
		return
	}
	// 使用磁盘文件
	file, err := os.Open(f.path)
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

// 显示目录列表
func (s *Server) listDir(r *Request, f http.File) {
	files, err := f.Readdir(-1)
	if err != nil {
		r.Response.WriteStatus(http.StatusInternalServerError, "Error reading directory")
		return
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
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
