package ghttp

import (
    "os"
    "fmt"
    "sort"
    "strings"
    "net/url"
    "net/http"
    "path/filepath"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/encoding/ghtml"
    "reflect"
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
    // 构造请求/返回参数对象
    request  := &ClientRequest{}
    response := &ServerResponse{}
    request.Request         = *r
    response.ResponseWriter = w
    if h := s.getHandler(gDEFAULT_DOMAIN, r.Method, r.URL.Path); h != nil {
        s.callHandler(h, request, response)
    } else {
        if h := s.getHandler(strings.Split(r.Host, ":")[0], r.Method, r.URL.Path); h != nil {
            s.callHandler(h, request, response)
        } else {
            s.serveFile(w, r)
        }
    }
}

// 初始化控制器
func (s *Server)callHandler(h *HandlerItem, r *ClientRequest, w *ServerResponse) {
    if h.faddr == nil {
        c := reflect.New(h.ctype)
        c.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(r), reflect.ValueOf(w)})
        c.MethodByName(h.fname).Call(nil)
        c.MethodByName("Shut").Call(nil)
    } else {
        h.faddr(s, r, w)
    }
    // 路由规则打包
    if buffer, err := s.Router.Patch(w.Buffer()); err == nil {
        w.ClearBuffer()
        w.Write(buffer)
    }
    // 输出缓冲区
    w.OutputBuffer()
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