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
func (h gtHttpServer)New(addr string, handler *http.Handler) error {
    return http.ListenAndServe(addr, handler)
}

// 创建一个HTTPserver
func (h gtHttpServer)NewByConfig(httpServerConfig *http.Server) error {
    return httpServerConfig.ListenAndServe()
}

// 绑定URI到操作函数/方法
func (h gtHttpServer)BindHandle(pattern string, handler *http.HandlerFunc ) {
    //http.HandleFunc(pattern, handler)
}


