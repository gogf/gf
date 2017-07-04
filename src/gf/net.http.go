package gf

import (
    "net/http"
    "time"
    "crypto/tls"
    "net"
    "log"
    "sync"
)

// 全局http封装对象
var Http gstHttp

// HTTP 结构体
type gstHttp struct {}

// HTTP Server 结构体
type gstHttpServer struct {
    // HTTP Server基础字段
    Addr            string
    Handler         http.Handler
    TLSConfig      *tls.Config
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    IdleTimeout     time.Duration
    MaxHeaderBytes  int
    ErrorLog       *log.Logger
    // gf 扩展信息字段
    ServerAgent     string
}

// 默认HTTP Server
var defaultHttpServer = gstHttpServer {
    Addr           : ":80",
    Handler        : http.HandlerFunc(defaultHttpHandle),
    ReadTimeout    : 10 * time.Second,
    WriteTimeout   : 10 * time.Second,
    IdleTimeout    : 10 * time.Second,
    MaxHeaderBytes : 1024,
    ServerAgent    : "gf",
}

// 默认HTTP Server处理入口
func defaultHttpHandle(w http.ResponseWriter, req *http.Request) {

}

// 获得一个默认的HTTP Server
func (h gstHttp)NewServer(addr string) (*gstHttpServer) {
    s      := defaultHttpServer
    s.Addr  = addr
    return &s
}



// 执行
func (h *gstHttpServer)Run(httpServerConfig *http.Server) error {
    return httpServerConfig.ListenAndServe()
}

// 绑定URI到操作函数/方法
func (h *gstHttpServer)BindHandle(pattern string, handler http.HandlerFunc ) {
    http.HandleFunc(pattern, handler)
}

// 通过映射数组绑定URI到操作函数/方法
func (h *gstHttpServer)BindHandleByMap(m map[string]http.HandlerFunc ) {
    for p, f := range m {
        h.BindHandle(p, f)
    }
}

