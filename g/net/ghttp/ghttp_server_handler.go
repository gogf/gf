<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 请求处理.
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package ghttp

import (
<<<<<<< HEAD
    "os"
    "fmt"
    "sort"
    "reflect"
    "strings"
    "net/url"
    "net/http"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/encoding/ghtml"
=======
    "fmt"
    "github.com/gogf/gf/g/encoding/ghtml"
    "github.com/gogf/gf/g/os/gspath"
    "github.com/gogf/gf/g/os/gtime"
    "net/http"
    "os"
    "reflect"
    "sort"
    "strings"
>>>>>>> upstream/master
)

// 默认HTTP Server处理入口，http包底层默认使用了gorutine异步处理请求，所以这里不再异步执行
func (s *Server)defaultHttpHandle(w http.ResponseWriter, r *http.Request) {
    s.handleRequest(w, r)
}

<<<<<<< HEAD
// 执行处理HTTP请求
=======
// 执行处理HTTP请求，
>>>>>>> upstream/master
// 首先，查找是否有对应域名的处理接口配置；
// 其次，如果没有对应的自定义处理接口配置，那么走默认的域名处理接口配置；
// 最后，如果以上都没有找到处理接口，那么进行文件处理；
func (s *Server)handleRequest(w http.ResponseWriter, r *http.Request) {
<<<<<<< HEAD
    // 去掉末尾的"/"号
    if r.URL.Path != "/" {
        r.URL.Path = strings.TrimRight(r.URL.Path, "/")
=======
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
        for r.URL.Path[len(r.URL.Path) - 1] == '/' {
            r.URL.Path = r.URL.Path[:len(r.URL.Path) - 1]
        }
>>>>>>> upstream/master
    }

    // 创建请求处理对象
    request := newRequest(s, r, w)

<<<<<<< HEAD
    // 错误日志使用recover进行判断
    defer func() {
        if request.LeaveTime == 0 {
            request.LeaveTime = gtime.Microsecond()
        }
        if e := recover(); e != nil {
            s.handleErrorLog(e, request)
        }
        s.handleAccessLog(request)
    }()

    // 事件 - BeforeServe
    s.callHookHandler(request, "BeforeServe")
    if h := s.getHandler(request); h != nil {
        s.callHandler(h, request)
    } else {
        s.serveFile(request)
    }
    // 事件 - AfterServe
    s.callHookHandler(request, "AfterServe")

    // 设置请求完成时间
    request.LeaveTime = gtime.Microsecond()

    // 事件 - BeforeOutput
    s.callHookHandler(request, "BeforeOutput")
    // 输出Cookie
    request.Cookie.Output()
    // 输出缓冲区
    request.Response.OutputBuffer()
    // 事件 - AfterOutput
    s.callHookHandler(request, "AfterOutput")

    // 将Request对象指针丢到队列中异步处理
    s.closeQueue.PushBack(request)
}

// 初始化控制器
func (s *Server)callHandler(h *HandlerItem, r *Request) {
    if h.faddr == nil {
        // 新建一个控制器对象处理请求
        c := reflect.New(h.ctype)
        c.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(r)})
        if !r.IsExited() {
            c.MethodByName(h.fname).Call(nil)
        }
        c.MethodByName("Shut").Call([]reflect.Value{reflect.ValueOf(r)})
    } else {
        if !r.IsExited() {
            h.faddr(r)
=======
    defer func() {
        // 设置请求完成时间
        request.LeaveTime = gtime.Microsecond()
        // 事件 - BeforeOutput
        if !request.IsExited() {
            s.callHookHandler(HOOK_BEFORE_OUTPUT, request)
        }
        // 如果没有产生异常状态，那么设置返回状态为200
        if request.Response.Status == 0 {
	        request.Response.Status = http.StatusOK
        }
        // error log
        if e := recover(); e != nil {
            request.Response.WriteStatus(http.StatusInternalServerError)
            s.handleErrorLog(e, request)
        }
        // access log
        s.handleAccessLog(request)
        // 输出Cookie
        request.Cookie.Output()
        // 输出缓冲区
        request.Response.Output()
        // 事件 - AfterOutput
        if !request.IsExited() {
            s.callHookHandler(HOOK_AFTER_OUTPUT, request)
        }
        // 更新Session会话超时时间
        request.Session.UpdateExpire()
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
    if !request.isFileRequest || isStaticDir {
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
        // 需要再次判断文件是否真实存在，
        // 因为文件检索可能使用了缓存，从健壮性考虑这里需要二次判断
        if request.isFileRequest /* && gfile.Exists(staticFile) */{
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
>>>>>>> upstream/master
        }
    }
}

<<<<<<< HEAD
// 处理静态文件请求
func (s *Server)serveFile(r *Request) {
    uri := r.URL.Path
    if s.config.ServerRoot != "" {
        // 获取文件的绝对路径
        path := strings.TrimRight(s.config.ServerRoot, gfile.Separator)
        if gfile.Separator != "/" {
            uri = strings.Replace(uri, "/", gfile.Separator, -1)
        }
        path = path + uri
        path = gfile.RealPath(path)
        if path != "" {
            // 文件/目录访问安全限制：服务的路径必须在ServerRoot下，否则会报错
            if len(path) >= len(s.config.ServerRoot) && strings.EqualFold(path[0 : len(s.config.ServerRoot)], s.config.ServerRoot) {
                s.doServeFile(r, path)
            } else {
                r.Response.WriteStatus(http.StatusForbidden)
            }
        } else {
            r.Response.WriteStatus(http.StatusNotFound)
        }
    } else {
        r.Response.WriteStatus(http.StatusNotFound)
    }
}

// http server静态文件处理
func (s *Server)doServeFile(r *Request, path string) {
    f, err := os.Open(path)
    if err != nil {
        return
    }
    info, _ := f.Stat()
    if info.IsDir() {
        if len(s.config.IndexFiles) > 0 {
            for _, file := range s.config.IndexFiles {
                fpath := path + gfile.Separator + file
                if gfile.Exists(fpath) {
                    f.Close()
                    s.doServeFile(r, fpath)
                    return
                }
            }
        }
=======
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
>>>>>>> upstream/master
        if s.config.IndexFolder {
            s.listDir(r, f)
        } else {
            r.Response.WriteStatus(http.StatusForbidden)
        }
    } else {
        // 读取文件内容返回, no buffer
<<<<<<< HEAD
        http.ServeContent(r.Response.Writer, &r.Request, info.Name(), info.ModTime(), f)
    }
    f.Close()
}

// 目录列表
func (s *Server)listDir(r *Request, f http.File) {
    dirs, err := f.Readdir(-1)
=======
        http.ServeContent(r.Response.Writer, r.Request, info.Name(), info.ModTime(), f)
    }
}

// 显示目录列表
func (s *Server)listDir(r *Request, f http.File) {
    files, err := f.Readdir(-1)
>>>>>>> upstream/master
    if err != nil {
        r.Response.WriteStatus(http.StatusInternalServerError, "Error reading directory")
        return
    }
<<<<<<< HEAD
    sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

    r.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
    r.Response.Write("<pre>\n")
    for _, d := range dirs {
        name := d.Name()
        if d.IsDir() {
            name += "/"
        }
        u := url.URL{Path: name}
        r.Response.Write(fmt.Sprintf("<a href=\"%s\">%s</a>\n", u.String(), ghtml.SpecialChars(name)))
    }
    r.Response.Write("</pre>\n")
}

// 开启异步队列处理循环，该异步线程与Server同生命周期
func (s *Server) startCloseQueueLoop() {
    go func() {
        for {
            if v := s.closeQueue.PopFront(); v != nil {
                r := v.(*Request)
                s.callHookHandler(r, "BeforeClose")
                // 关闭当前会话的Cookie
                r.Cookie.Close()
                // 更新Session会话超时时间
                r.Session.UpdateExpire()
                s.callHookHandler(r, "AfterClose")
            }
        }
    }()
}
=======
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
>>>>>>> upstream/master
