package ghttp

import (
    "net/http"
    "strings"
    "path/filepath"
    "g/os/gfile"
    "g/core/gutil"
    "os"
    "fmt"
    "sort"
    "net/url"
)

// 默认HTTP Server处理入口，底层默认使用了gorutine调用该接口
func (s *Server)defaultHttpHandle(w http.ResponseWriter, r *http.Request) {
    if f, ok := s.handlerMap[r.URL.String()]; ok {
        f(w, r)
    } else {
        s.serveFile(w, r)
    }
}

// 处理静态文件请求
func (s *Server)serveFile(w http.ResponseWriter, r *http.Request) {
    uri := r.URL.String()
    if s.config.ServerRoot != "" {
        // 获取文件的绝对路径
        path := strings.TrimRight(s.config.ServerRoot, string(filepath.Separator))
        path  = path + uri
        path  = gfile.RealPath(path)
        if (path != "") {
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
        s.listDir(w, f)
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
        url := url.URL{Path: name}
        fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), gutil.HtmlSpecialChars(name))
    }
    fmt.Fprintf(w, "</pre>\n")
}

// 404
func (s *Server)NotFound(w http.ResponseWriter, r *http.Request) {
    http.NotFound(w, r)
}