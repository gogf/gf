package ghttp

import (
    "net/http"
    "strings"
    "path/filepath"
    "crypto/tls"
    "time"
    "log"
    "sync"
    "errors"
    "gitee.com/johng/gf/g/container/gmap"
)

// http server结构体
type Server struct {
    hmu        sync.RWMutex // handlerMap互斥锁
    name       string       // 服务名称，方便识别
    server     http.Server  // 底层http server对象
    config     ServerConfig // 配置对象
    handlerMap HandlerMap   // 回调函数
    status     int8         // 当前服务器状态(0：未启动，1：运行中)
}

// http回调函数
type HandlerFunc func(*ClientRequest, *ServerResponse)

// uri与回调函数的绑定记录表
type HandlerMap map[string]HandlerFunc

// Server表，用以存储和检索名称与Server对象之间的关联关系
var serverMapping *gmap.StringInterfaceMap = gmap.NewStringInterfaceMap()

// 创建一个默认配置的HTTP Server(默认监听端口是80)
func GetServer(name string) (*Server) {
    if s := serverMapping.Get(name); s != nil {
        return s.(*Server)
    }
    s     := &Server{}
    s.name = name
    s.SetConfig(defaultServerConfig)
    serverMapping.Set(name, s)
    return s
}

// 执行
func (s *Server) Run() error {
    if s.status == 1 {
        return errors.New("server is already running")
    }

    // 底层http server配置
    if s.config.Handler == nil {
        s.config.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }
    // 底层http server初始化
    s.server  = http.Server {
        Addr           : s.config.Addr,
        Handler        : s.config.Handler,
        ReadTimeout    : s.config.ReadTimeout,
        WriteTimeout   : s.config.WriteTimeout,
        IdleTimeout    : s.config.IdleTimeout,
        MaxHeaderBytes : s.config.MaxHeaderBytes,
    }
    // 执行端口监听
    if err := s.server.ListenAndServe(); err != nil {
        return err
    }
    s.status = 1
    return nil
}

// http server setting设置
// 注意使用该方法进行http server配置时，需要配置所有的配置项，否则没有配置的属性将会默认变量为空
func (s *Server)SetConfig(c ServerConfig) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
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
    return nil
}

// 设置http server参数 - Addr
func (s *Server)SetAddr(addr string) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.Addr = addr
    return nil
}

// 设置http server参数 - TLSConfig
func (s *Server)SetTLSConfig(tls *tls.Config) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.TLSConfig = tls
    return nil
}

// 设置http server参数 - ReadTimeout
func (s *Server)SetReadTimeout(t time.Duration) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.ReadTimeout = t
    return nil
}

// 设置http server参数 - WriteTimeout
func (s *Server)SetWriteTimeout(t time.Duration) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.WriteTimeout = t
    return nil
}

// 设置http server参数 - IdleTimeout
func (s *Server)SetIdleTimeout(t time.Duration) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.IdleTimeout = t
    return nil
}

// 设置http server参数 - MaxHeaderBytes
func (s *Server)SetMaxHeaderBytes(b int) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.MaxHeaderBytes = b
    return nil
}

// 设置http server参数 - ErrorLog
func (s *Server)SetErrorLog(logger *log.Logger) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.ErrorLog = logger
    return nil
}

// 设置http server参数 - IndexFiles
func (s *Server)SetIndexFiles(index []string) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.IndexFiles = index
    return nil
}

// 设置http server参数 - IndexFolder
func (s *Server)SetIndexFolder(index bool) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.IndexFolder = index
    return nil
}

// 设置http server参数 - ServerAgent
func (s *Server)SetServerAgent(agent string) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.ServerAgent = agent
    return nil
}

// 设置http server参数 - ServerRoot
func (s *Server)SetServerRoot(root string) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.ServerRoot  = strings.TrimRight(root, string(filepath.Separator))
    return nil
}

// 绑定URI到操作函数/方法
// pattern的格式形如：/user/list, put:/user, delete:/user
// 支持RESTful的请求格式，具体业务逻辑由绑定的处理方法来执行
func (s *Server)BindHandle(pattern string, handler HandlerFunc) error {
    if s.status == 1 {
        return errors.New("server handlers cannot be changed while running")
    }

    s.hmu.Lock()
    defer s.hmu.Unlock()

    if s.handlerMap == nil {
        s.handlerMap = make(HandlerMap)
    }
    key    := ""
    result := strings.Split(pattern, ":")
    if len(result) > 1 {
        key = strings.ToUpper(result[0]) + ":" + result[1]
    } else {
        key = strings.TrimSpace(pattern)
    }
    if _, ok := s.handlerMap[key]; ok {
        return errors.New("duplicated http server handler for: " + pattern)
    } else {
        s.handlerMap[key] = handler
    }
    return nil
}

// 通过映射数组绑定URI到操作函数/方法
func (s *Server)BindHandleByMap(m HandlerMap) {
    for p, f := range m {
        s.BindHandle(p, f)
    }
}

// 绑定控制器，控制器需要继承gmvc.ControllerBase对象并实现需要的REST方法
func (s *Server)BindControllerRest(uri string, c ControllerRest) {
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


