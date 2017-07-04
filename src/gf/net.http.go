package gf

import (
    "net/http"
)

// HTTP Server 结构体
type gstHttpServer struct {}

// HTTP 结构体
type gstHttp struct {
    Server gstHttpServer
}

// 全局http封装对象
var Http gstHttp

// 创建一个简单的HTTPserver
func (h gstHttpServer)Start(addr string) error {
    return http.ListenAndServe(addr, nil)
}

// 创建一个自定义配置的HTTPserver
// 常用的配置包括：
// Addr,
// ReadTimeout,
// WriteTimeout,
// IdleTimeout,
// MaxHeaderBytes,
// ErrorLog
func (h gstHttpServer)StartByConfig(httpServerConfig *http.Server) error {
    return httpServerConfig.ListenAndServe()
}

// 绑定URI到操作函数/方法
func (h gstHttpServer)BindHandle(pattern string, handler http.HandlerFunc ) {
    http.HandleFunc(pattern, handler)
}

// 通过映射数组绑定URI到操作函数/方法
func (h gstHttpServer)BindHandleByMap(m map[string]http.HandlerFunc ) {
    for p, f := range m {
        h.BindHandle(p, f)
    }
}

