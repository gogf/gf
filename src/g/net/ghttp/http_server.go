package ghttp

import (
    "net/http"
    "strings"
    "path/filepath"
)

// 执行
func (h *Server)Run() error {
    // 底层http server配置
    h.server  = http.Server {
        Addr           : h.config.Addr,
        Handler        : h.config.Handler,
        ReadTimeout    : h.config.ReadTimeout,
        WriteTimeout   : h.config.WriteTimeout,
        IdleTimeout    : h.config.IdleTimeout,
        MaxHeaderBytes : h.config.MaxHeaderBytes,
    }
    // 执行端口监听
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
    h.config = c
    if h.config.ServerRoot != "" {
        h.SetServerRoot(h.config.ServerRoot)
    }
}

// 设置http server参数 - IndexFiles
func (h *Server)SetIndexFiles(index []string) {
    h.config.IndexFiles = index
}

// 设置http server参数 - IndexFolder
func (h *Server)SetIndexFolder(index bool) {
    h.config.IndexFolder = index
}

// 设置http server参数 - ServerAgent
func (h *Server)SetServerAgent(agent string) {
    h.config.ServerAgent = agent
}

// 设置http server参数 - ServerRoot
func (h *Server)SetServerRoot(root string) {
    h.config.ServerRoot  = strings.TrimRight(root, string(filepath.Separator))
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

