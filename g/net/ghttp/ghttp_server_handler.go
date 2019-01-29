// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 请求处理.

package ghttp

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/ghtml"
    "gitee.com/johng/gf/g/os/gspath"
    "gitee.com/johng/gf/g/os/gtime"
    "net/http"
    "os"
    "reflect"
    "sort"
    "strings"
)

// 默认HTTP Server处理入口，http包底层默认使用了gorutine异步处理请求，所以这里不再异步执行
func (s *Server)defaultHttpHandle(w http.ResponseWriter, r *http.Request) {
    s.handleRequest(w, r)
}

// 执行处理HTTP请求
// 首先，查找是否有对应域名的处理接口配置；
// 其次，如果没有对应的自定义处理接口配置，那么走默认的域名处理接口配置；
// 最后，如果以上都没有找到处理接口，那么进行文件处理；
func (s *Server)handleRequest(w http.ResponseWriter, r *http.Request) {
    // 重写规则判断
    if len(s.config.Rewrites) > 0 {
        if rewrite, ok := s.config.Rewrites[r.URL.Path]; ok {
            r.URL.Path = rewrite
        }
    }
    // 去掉末尾的"/"号
    if r.URL.Path != "/" {
        for r.URL.Path[len(r.URL.Path) - 1] == '/' {
            r.URL.Path = r.URL.Path[:len(r.URL.Path) - 1]
        }
    }

    // 创建请求处理对象
    request := newRequest(s, r, w)

    defer func() {
        if request.LeaveTime == 0 {
            request.LeaveTime = gtime.Microsecond()
        }
        s.callHookHandler(HOOK_BEFORE_CLOSE, request)
        // access log
        s.handleAccessLog(request)
        // error log使用recover进行判断
        if e := recover(); e != nil {
            s.handleErrorLog(e, request)
        }
        // 更新Session会话超时时间
        request.Session.UpdateExpire()
        s.callHookHandler(HOOK_AFTER_CLOSE, request)
    }()

    // ============================================================
    // 优先级控制:
    // 静态文件 > 动态服务 > 静态目录
    // ============================================================

    staticFile  := ""
    isStaticDir := false
    // 优先执行静态文件检索(检测是否存在对应的静态文件，包括index files处理)
    if s.config.FileServerEnabled {
        staticFile, isStaticDir = s.searchStaticFile(r.URL.Path)
        if staticFile != "" {
            request.isFileRequest = true
        }
    }

    // 动态服务检索
    handler := (*handlerItem)(nil)
    if !request.IsFileRequest() || isStaticDir {
        if parsedItem := s.getServeHandlerWithCache(request); parsedItem != nil {
            handler = parsedItem.handler
            for k, v := range parsedItem.values {
                request.routerVars[k] = v
            }
            request.Router = parsedItem.handler.router
        }
    }

    // 判断最终对该请求提供的服务方式
    if isStaticDir && handler != nil {
        request.isFileRequest = false
    }

    // 事件 - BeforeServe
    s.callHookHandler(HOOK_BEFORE_SERVE, request)

    // 执行静态文件服务/回调控制器/执行对象/方法
    if !request.IsExited() {
        // 需要再次判断文件是否真实存在，因为文件检索可能使用了缓存，从健壮性考虑这里需要二次判断
        if request.isFileRequest /* && gfile.Exists(staticFile) */{
            // 静态文件
            s.serveFile(request, staticFile)
        } else {
            if handler != nil {
                // 动态服务
                s.callServeHandler(handler, request)
            } else {
                if isStaticDir {
                    // 静态目录
                    s.serveFile(request, staticFile)
                } else {
                    if len(request.Response.Header()) == 0 &&
                        request.Response.Status == 0 &&
                        request.Response.BufferLength() == 0 {
                        request.Response.WriteStatus(http.StatusNotFound)
                    }
                }
            }
        }
    }

    // 事件 - AfterServe
    if !request.IsExited() {
        s.callHookHandler(HOOK_AFTER_SERVE, request)
    }

    // 设置请求完成时间
    request.LeaveTime = gtime.Microsecond()

    // 事件 - BeforeOutput
    if !request.IsExited() {
        s.callHookHandler(HOOK_BEFORE_OUTPUT, request)
    }
    // 输出Cookie
    request.Cookie.Output()
    // 输出缓冲区
    request.Response.OutputBuffer()
    // 事件 - AfterOutput
    if !request.IsExited() {
        s.callHookHandler(HOOK_AFTER_OUTPUT, request)
    }
}

// 查找静态文件的绝对路径
func (s *Server) searchStaticFile(uri string) (filePath string, isDir bool) {
    // 优先查找URI映射
    if len(s.config.StaticPaths) > 0 {
        for _, item := range s.config.StaticPaths {
            if len(uri) >= len(item.prefix) && strings.EqualFold(item.prefix, uri[0 : len(item.prefix)]) {
                // 防止类似 /static/style 映射到 /static/style.css 的情况
                if len(uri) > len(item.prefix) && uri[len(item.prefix)] != '/' {
                    continue
                }
                return gspath.Search(item.path, uri[len(item.prefix):], s.config.IndexFiles...)
            }
        }
    }
    // 其次查找root和search path
    if len(s.config.SearchPaths) > 0 {
        for _, path := range s.config.SearchPaths {
            if filePath, isDir = gspath.Search(path, uri, s.config.IndexFiles...); filePath != "" {
                return filePath, isDir
            }
        }
    }
    return "", false
}

// 调用服务接口
func (s *Server) callServeHandler(h *handlerItem, r *Request) {
    if h.faddr == nil {
        c := reflect.New(h.ctype)
        s.niceCallFunc(func() {
            c.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(r)})
        })
        if !r.IsExited() {
            s.niceCallFunc(func() {
                c.MethodByName(h.fname).Call(nil)
            })
        }
        if !r.IsExited() {
            s.niceCallFunc(func() {
                c.MethodByName("Shut").Call(nil)
            })
        }
    } else {
        if h.finit != nil {
            s.niceCallFunc(func() {
                h.finit(r)
            })
        }
        if !r.IsExited() {
            s.niceCallFunc(func() {
                h.faddr(r)
            })
        }
        if h.fshut != nil && !r.IsExited() {
            s.niceCallFunc(func() {
                h.fshut(r)
            })
        }
    }
}

// 友好地调用方法
func (s *Server) niceCallFunc(f func()) {
    defer func() {
        if err := recover(); err != nil {
            switch err {
                case gEXCEPTION_EXIT:    fallthrough
                case gEXCEPTION_EXIT_ALL:
                    return
                default:
                    panic(err)
            }
        }
    }()
    f()
}

// http server静态文件处理，path可以为相对路径也可以为绝对路径
func (s *Server) serveFile(r *Request, path string) {
    f, err := os.Open(path)
    if err != nil {
        r.Response.WriteStatus(http.StatusForbidden)
        return
    }
    defer f.Close()
    info, _ := f.Stat()
    if info.IsDir() {
        if s.config.IndexFolder {
            s.listDir(r, f)
        } else {
            r.Response.WriteStatus(http.StatusForbidden)
        }
    } else {
        // 读取文件内容返回, no buffer
        http.ServeContent(r.Response.Writer, &r.Request, info.Name(), info.ModTime(), f)
    }
}

// 显示目录列表
func (s *Server)listDir(r *Request, f http.File) {
    files, err := f.Readdir(-1)
    if err != nil {
        r.Response.WriteStatus(http.StatusInternalServerError, "Error reading directory")
        return
    }
    sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

    r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
    r.Response.Write("<pre>\n")
    if r.URL.Path != "/" {
        r.Response.Write(fmt.Sprint("<a href=\"..\">..</a>\n"))
    }
    for _, file := range files {
        name := file.Name()
        if file.IsDir() {
            name += "/"
        }
        r.Response.Write(fmt.Sprintf("<a href=\"%s\">%s</a>\n", name, ghtml.SpecialChars(name)))
    }
    r.Response.Write("</pre>\n")
}
