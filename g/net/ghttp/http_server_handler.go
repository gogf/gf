// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package ghttp

import (
    "os"
    "fmt"
    "sort"
    "reflect"
    "strings"
    "net/url"
    "net/http"
    "path/filepath"
    "gitee.com/johng/gf/g/os/gfile"
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
    // 路由解析
    uri         := r.URL.String()
    result, err := s.Router.Dispatch(uri)
    if err == nil && strings.Compare(uri, result) != 0 {
        r.URL, _ = r.URL.Parse(result)
    }
    // 构造请求参数对象
    request  := &Request{
        Id       : s.idgen.Int(),
        Server   : s,
        Request  : *r,
        Response : &Response {
            ResponseWriter : w,
        },
    }
    if h := s.getHandler(gDEFAULT_DOMAIN, r.Method, r.URL.Path); h != nil {
        s.callHandler(h, request)
    } else {
        if h := s.getHandler(strings.Split(r.Host, ":")[0], r.Method, r.URL.Path); h != nil {
            s.callHandler(h, request)
        } else {
            s.serveFile(w, r)
        }
    }
}

// 初始化控制器
func (s *Server)callHandler(h *HandlerItem, r *Request) {
    // 会话处理
    r.Cookie  = NewCookie(r)
    r.Session = GetSession(r.Cookie.SessionId())

    // 请求处理
    if h.faddr == nil {
        // 新建一个控制器对象处理请求
        c := reflect.New(h.ctype)
        c.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(r)})
        c.MethodByName(h.fname).Call(nil)
        c.MethodByName("Shut").Call(nil)
    } else {
        // 直接调用注册的方法处理请求
        h.faddr(r)
    }
    // 路由规则打包
    if buffer, err := s.Router.Patch(r.Response.Buffer()); err == nil {
        r.Response.ClearBuffer()
        r.Response.Write(buffer)
    }

    // 输出Cookie
    r.Cookie.Output()

    // 输出缓冲区
    r.Response.OutputBuffer()

    // 将Request对象指针丢到队列中异步处理
    s.closeQueue.PushBack(r)
}

// 处理静态文件请求
func (s *Server)serveFile(w http.ResponseWriter, r *http.Request) {
    uri := r.URL.String()
    if s.config.ServerRoot != "" {
        // 获取文件的绝对路径
        path := strings.TrimRight(s.config.ServerRoot, string(filepath.Separator))
        path  = path + uri
        path  = gfile.RealPath(path)
        if path != "" {
            s.doServeFile(w, r, path)
        } else {
            s.NotFound(w, r)
        }
    } else {
        s.NotFound(w, r)
    }
}

// http server静态文件处理
func (s *Server)doServeFile(w http.ResponseWriter, r *http.Request, path string) {
    f, err := os.Open(path)
    if err != nil {
        return
    }
    info, _ := f.Stat()
    if info.IsDir() {
        if len(s.config.IndexFiles) > 0 {
            for _, file := range s.config.IndexFiles {
                fpath := path + "/" + file
                if gfile.Exists(fpath) {
                    f.Close()
                    s.doServeFile(w, r, fpath)
                    return
                }
            }
        }
        if s.config.IndexFolder {
            s.listDir(w, f)
        } else {
            s.ResponseStatus(w, http.StatusForbidden)
        }
    } else {
        http.ServeContent(w, r, info.Name(), info.ModTime(), f)
    }
    f.Close()
}

// 目录列表
func (s *Server)listDir(w http.ResponseWriter, f http.File) {
    dirs, err := f.Readdir(-1)
    if err != nil {
        http.Error(w, "Error reading directory", http.StatusInternalServerError)
        return
    }
    sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, "<pre>\n")
    for _, d := range dirs {
        name := d.Name()
        if d.IsDir() {
            name += "/"
        }
        u := url.URL{Path: name}
        fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", u.String(), ghtml.SpecialChars(name))
    }
    fmt.Fprintf(w, "</pre>\n")
}

// 返回http状态码，并使用默认配置的字符串返回信息
func (s *Server)ResponseStatus(w http.ResponseWriter, code int) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.WriteHeader(code)
    w.Write([]byte(http.StatusText(code)))
}

// 404
func (s *Server)NotFound(w http.ResponseWriter, r *http.Request) {
    http.NotFound(w, r)
}