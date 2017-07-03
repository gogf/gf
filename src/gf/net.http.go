package gf

import (
    "net/http"
)

// HTTP Server 结构体
type gtHttpServer struct {}

// HTTP 结构体
type gtHttp struct {
    Server gtHttpServer
}

// 全局http封装对象
var Http gtHttp

// 创建一个简单的HTTPserver
func (h gtHttpServer)Start(addr string) error {
    return http.ListenAndServe(addr, nil)
}

// 创建一个自定义配置的HTTPserver
func (h gtHttpServer)StartByConfig(httpServerConfig *http.Server) error {
    return httpServerConfig.ListenAndServe()
}

// 绑定URI到操作函数/方法
func (h gtHttpServer)BindHandle(pattern string, handler http.HandlerFunc ) {
    http.HandleFunc(pattern, handler)
}


