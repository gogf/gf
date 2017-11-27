package ghttp

import (
    "net/http"
    "strings"
    "path/filepath"
    "crypto/tls"
    "time"
    "log"
    "regexp"
    "gitee.com/johng/gf/g/os/glog"
)

// 执行
func (s *Server)Run() error {
    // 底层http server配置
    if s.config.Handler == nil {
        s.config.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }
    s.server  = http.Server {
        Addr           : s.config.Addr,
        Handler        : s.config.Handler,
        ReadTimeout    : s.config.ReadTimeout,
        WriteTimeout   : s.config.WriteTimeout,
        IdleTimeout    : s.config.IdleTimeout,
        MaxHeaderBytes : s.config.MaxHeaderBytes,
    }
    // 执行端口监听
    err := s.server.ListenAndServe()
    if err != nil {
        glog.Fatalln(err)
    }
    return err
}

// 获取默认的http server设置
func (h Server)GetDefaultSetting() ServerConfig {
    return defaultServerConfig
}

// http server setting设置
// 注意使用该方法进行http server配置时，需要配置所有的配置项，否则没有配置的属性将会默认变量为空
func (s *Server)SetConfig(c ServerConfig) {
    if c.Handler == nil {
        c.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }
    s.config = c
    // 需要处理server root最后的目录分隔符号
    if s.config.ServerRoot != "" {
        s.SetServerRoot(s.config.ServerRoot)
    }
    // 必需设置默认值的属性
    if len(s.config.IndexFiles) < 1 {
        s.SetIndexFiles(defaultServerConfig.IndexFiles)
    }
    if s.config.ServerAgent == "" {
        s.SetServerAgent(defaultServerConfig.ServerAgent)
    }
}

// 设置http server参数 - Addr
func (s *Server)SetAddr(addr string) {
    s.config.Addr = addr
}

// 设置http server参数 - Handler
func (s *Server)SetHandler(handler http.Handler) {
    s.config.Handler = handler
}

// 设置http server参数 - TLSConfig
func (s *Server)SetTLSConfig(tls *tls.Config) {
    s.config.TLSConfig = tls
}

// 设置http server参数 - ReadTimeout
func (s *Server)SetReadTimeout(t time.Duration) {
    s.config.ReadTimeout = t
}

// 设置http server参数 - WriteTimeout
func (s *Server)SetWriteTimeout(t time.Duration) {
    s.config.WriteTimeout = t
}

// 设置http server参数 - IdleTimeout
func (s *Server)SetIdleTimeout(t time.Duration) {
    s.config.IdleTimeout = t
}

// 设置http server参数 - MaxHeaderBytes
func (s *Server)SetMaxHeaderBytes(b int) {
    s.config.MaxHeaderBytes = b
}

// 设置http server参数 - ErrorLog
func (s *Server)SetErrorLog(logger *log.Logger) {
    s.config.ErrorLog = logger
}

// 设置http server参数 - IndexFiles
func (s *Server)SetIndexFiles(index []string) {
    s.config.IndexFiles = index
}

// 设置http server参数 - IndexFolder
func (s *Server)SetIndexFolder(index bool) {
    s.config.IndexFolder = index
}

// 设置http server参数 - ServerAgent
func (s *Server)SetServerAgent(agent string) {
    s.config.ServerAgent = agent
}

// 设置http server参数 - ServerRoot
func (s *Server)SetServerRoot(root string) {
    s.config.ServerRoot  = strings.TrimRight(root, string(filepath.Separator))
}

// 绑定URI到操作函数/方法
// pattern的格式形如：/user/list, put:/user, delete:/user
// 支持RESTful的请求格式，具体业务逻辑由绑定的处理方法来执行
func (s *Server)BindHandle(pattern string, handler HandlerFunc )  {
    if s.handlerMap == nil {
        s.handlerMap = make(HandlerMap)
    }
    key    := ""
    reg    := regexp.MustCompile(`(\w+?)\s*:\s*(.+)`)
    result := reg.FindStringSubmatch(pattern)
    if len(result) > 1 {
        key = strings.ToUpper(result[1]) + ":" + result[2]
    } else {
        key = strings.TrimSpace(pattern)
    }
    if _, ok := s.handlerMap[key]; ok {
        panic("duplicated http server handler for: " + pattern)
    } else {
        s.handlerMap[key] = handler
    }
}

// 通过映射数组绑定URI到操作函数/方法
func (s *Server)BindHandleByMap(m HandlerMap) {
    for p, f := range m {
        s.BindHandle(p, f)
    }
}

// 绑定控制器，控制器需要继承gmvc.ControllerBase对象并实现需要的REST方法
func (s *Server)BindController(uri string, c ControllerApi) {
    s.BindHandleByMap(HandlerMap{
        "GET:"     + uri : c.Get,
        "PUT:"     + uri : c.Put,
        "POST:"    + uri : c.Post,
        "DELETE:"  + uri : c.Delete,
        "PATCH:"   + uri : c.Patch,
        "HEAD:"    + uri : c.Head,
        "CONNECT:" + uri : c.Connect,
        "OPTIONS:" + uri : c.Options,
        "TRACE:"   + uri : c.Trace,
    })
}


