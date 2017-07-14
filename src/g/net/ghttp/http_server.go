package ghttp

import (
    "net/http"
)

// 执行
func (h *Server)Run() error {
    err := h.server.ListenAndServe()
    if err != nil {
        panic(err)
    }
    return err
}

// 获取默认的http server设置
func (h Server)GetDefaultSetting() ServerConfig {
    return defaultServerConfig
}

// http server setting设置
func (h *Server)SetConfig(c ServerConfig) {
    if c.Handler == nil {
        c.Handler = http.HandlerFunc(h.defaultHttpHandle)
    }
    h.config  = c
    h.server  = http.Server {
        Addr           : c.Addr,
        Handler        : c.Handler,
        ReadTimeout    : c.ReadTimeout,
        WriteTimeout   : c.WriteTimeout,
        IdleTimeout    : c.IdleTimeout,
        MaxHeaderBytes : c.MaxHeaderBytes,
    }
}

// 设置http server参数
func (h *Server)SetServerAgent(agent string) {
    h.config.ServerAgent = agent
}

// 设置http server参数
func (h *Server)SetServerRoot(root string) {
    h.config.ServerRoot = root
}

// 绑定URI到操作函数/方法
func (h *Server)BindHandle(pattern string, handler http.HandlerFunc )  {
    if h.handlerMap == nil {
        h.handlerMap = make(map[string]http.HandlerFunc)
    }
    if _, ok := h.handlerMap[pattern]; ok {
        panic("duplicated http server handler for: " + pattern)
    } else {
        h.handlerMap[pattern] = handler
    }
}

// 通过映射数组绑定URI到操作函数/方法
func (h *Server)BindHandleByMap(m map[string]http.HandlerFunc ) {
    for p, f := range m {
        h.BindHandle(p, f)
    }
}

