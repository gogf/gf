package ghttp

import (
    "net/http"
    "strings"
    "path/filepath"
    "g/os/gfile"
)


// 默认HTTP Server处理入口，底层默认使用了gorutine调用该接口
func (h *Server)defaultHttpHandle(w http.ResponseWriter, r *http.Request) {
    if f, ok := h.handlerMap[r.URL.String()]; ok {
        f(w, r)
    } else {
        h.serveFile(w, r)
    }
}

// 处理静态文件请求
func (h *Server)serveFile(w http.ResponseWriter, r *http.Request) {
    uri := r.URL.String()
    if h.config.ServerRoot != "" {
        path := strings.TrimRight(h.config.ServerRoot, string(filepath.Separator))
        path  = path + uri
        // fmt.Println(path)
        if (gfile.Exists(path)) {
            http.ServeFile(w, r, path)
        } else {
            http.NotFound(w, r)
        }
    } else {
        panic("http server root is empty while handling static file request")
    }
}