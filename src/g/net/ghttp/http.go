package ghttp

import (
    "net/http"
    "time"
    "crypto/tls"
    "log"
    //"fmt"
    "strings"
    "path/filepath"
    "g/os/gfile"
)

// @todo 标准库对静态文件的处理性能比Nginx稍弱，不能使用标准库方法，需自行处理

// http server结构体
type Server struct {
    server     http.Server
    config     ServerConfig
    handlerMap map[string]http.HandlerFunc
}

// HTTP Server 设置结构体
type ServerConfig struct {
    // HTTP Server基础字段
    Addr            string        // 监听IP和端口，监听本地所有IP使用":端口"
    Handler         http.Handler  // 默认的处理函数
    TLSConfig      *tls.Config    // TLS配置
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    IdleTimeout     time.Duration
    MaxHeaderBytes  int           // 最大的header长度
    ErrorLog       *log.Logger    // 错误日志的处理接口
    // gf 扩展信息字段
    IndexFolder     bool          // 如果访问目录是否显示目录列表
    ServerAgent     string        // server agent
    ServerRoot      string        // 服务器服务的本地目录根路径
}

// 默认HTTP Server
var defaultServerConfig = ServerConfig {
    Addr           : ":80",
    Handler        : nil,
    ReadTimeout    : 10 * time.Second,
    WriteTimeout   : 10 * time.Second,
    IdleTimeout    : 10 * time.Second,
    MaxHeaderBytes : 1024,
    ServerAgent    : "gf",
    ServerRoot     : "",
}

// 修改默认的http server配置
func SetDefaultServerConfig (c ServerConfig) {
    defaultServerConfig = c
}

// 创建一个默认配置的HTTP Server(默认监听端口是80)
func New() (*Server) {
    return NewByConfig(defaultServerConfig)
}

// 创建一个HTTP Server，返回指针
func NewByAddr(addr string) (*Server) {
    config     := defaultServerConfig
    config.Addr = addr
    return NewByConfig(config)
}

// 创建一个HTTP Server
func NewByAddrRoot(addr string, root string) (*Server) {
    config           := defaultServerConfig
    config.Addr       = addr
    config.ServerRoot = root
    return NewByConfig(config)
}

// 根据输入配置创建一个http server对象
func NewByConfig(s ServerConfig) (*Server) {
    var server Server
    server.SetConfig(s)
    return &server
}

// 执行
func (h *Server)Run() error {
    err := h.server.ListenAndServe()
    if err != nil {
        panic(err)
    }
    return err
}

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

