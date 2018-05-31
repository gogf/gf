// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 请求处理.

package ghttp

import (
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
    // 去掉末尾的"/"号
    if r.URL.Path != "/" {
        r.URL.Path = strings.TrimRight(r.URL.Path, "/")
    }

    // 创建请求处理对象
    request := newRequest(s, r, w)

    defer func() {
        if request.LeaveTime == 0 {
            request.LeaveTime = gtime.Microsecond()
        }
        // access log
        s.handleAccessLog(request)
        // error log使用recover进行判断
        if e := recover(); e != nil {
            s.handleErrorLog(e, request)
        }
        // 将Request对象指针丢到队列中异步关闭
        s.closeQueue.PushBack(request)
    }()

    // 路由注册检索
    handler := s.getHandler(request)
    if handler == nil {
        // 如果路由不匹配，那么执行静态文件检索
        path := s.paths.Search(r.URL.Path)
        if path != "" {
            s.serveFile(request, path)
        } else {
            request.Response.WriteStatus(http.StatusNotFound)
            request.Response.OutputBuffer()
        }
        return
    }

    // **********************************************
    // 以下操作仅对路由控制有效，包括事件处理，不对静态文件有效
    // **********************************************

    // 事件 - BeforeServe
    s.callHookHandler(request, "BeforeServe")

    // 执行回调控制器/执行对象/方法
    s.callHandler(handler, request)

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
        }
    }
}

// http server静态文件处理
func (s *Server)serveFile(r *Request, path string) {
    f, err := os.Open(path)
    if err != nil {
        return
    }
    info, _ := f.Stat()
    if info.IsDir() {
        // 处理访问目录
        if len(s.config.IndexFiles) > 0 {
            for _, file := range s.config.IndexFiles {
                fpath := path + gfile.Separator + file
                if gfile.Exists(fpath) {
                    f.Close()
                    s.serveFile(r, fpath)
                    return
                }
            }
        }
        if s.config.IndexFolder {
            s.listDir(r, f)
        } else {
            r.Response.WriteStatus(http.StatusForbidden)
        }
    } else {
        // 读取文件内容返回, no buffer
        http.ServeContent(r.Response.Writer, &r.Request, info.Name(), info.ModTime(), f)
    }
    f.Close()
}

// 目录列表
func (s *Server)listDir(r *Request, f http.File) {
    dirs, err := f.Readdir(-1)
    if err != nil {
        r.Response.WriteStatus(http.StatusInternalServerError, "Error reading directory")
        return
    }
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