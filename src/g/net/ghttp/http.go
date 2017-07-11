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

// @todo 静态文件的处理性能比Nginx稍弱，不能使用标准库方法，需自行处理

// http server结构体
type HttpServer struct {
    server     http.Server
    setting    HttpServerSetting
    handlerMap map[string]http.HandlerFunc
}

// HTTP Server 设置结构体
type HttpServerSetting struct {
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
    ServerRoot      string
}

// 默认HTTP Server
var defaultHttpServerSetting = HttpServerSetting {
    Addr           : ":80",
    Handler        : nil,
    ReadTimeout    : 10 * time.Second,
    WriteTimeout   : 10 * time.Second,
    IdleTimeout    : 10 * time.Second,
    MaxHeaderBytes : 1024,
    ServerAgent    : "gf",
    ServerRoot     : "",
}


// 创建一个默认配置的HTTP Server(默认监听端口是80)
func New() (*HttpServer) {
    return NewBySetting(defaultHttpServerSetting)
}

// 创建一个HTTP Server，返回指针
func NewByAddr(addr string) (*HttpServer) {
    setting     := defaultHttpServerSetting
    setting.Addr = addr
    return NewBySetting(setting)
}

// 创建一个HTTP Server
func NewByAddrRoot(addr string, root string) (*HttpServer) {
    setting           := defaultHttpServerSetting
    setting.Addr       = addr
    setting.ServerRoot = root
    return NewBySetting(setting)
}

// 根据输入配置创建一个http server对象
func NewBySetting(s HttpServerSetting) (*HttpServer) {
    var server HttpServer
    server.SetSetting(s)
    return &server
}

// 执行
func (h *HttpServer)Run() error {
    err := h.server.ListenAndServe()
    if err != nil {
        panic(err)
    }
    return err
}

// 默认HTTP Server处理入口
func (h *HttpServer)defaultHttpHandle(w http.ResponseWriter, r *http.Request) {
    if f, ok := h.handlerMap[r.URL.String()]; ok {
        f(w, r)
    } else {
        h.serveFile(w, r)
    }
}

// 处理静态文件请求
func (h *HttpServer)serveFile(w http.ResponseWriter, r *http.Request) {
    uri := r.URL.String()
    if h.setting.ServerRoot != "" {
        path := strings.TrimRight(h.setting.ServerRoot, string(filepath.Separator))
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
func (h HttpServer)GetDefaultSetting() HttpServerSetting {
    return defaultHttpServerSetting
}

// http server setting设置
func (h *HttpServer)SetSetting(s HttpServerSetting) {
    if s.Handler == nil {
        s.Handler = http.HandlerFunc(h.defaultHttpHandle)
    }
    h.setting = s
    h.server  = http.Server {
        Addr           : s.Addr,
        Handler        : s.Handler,
        ReadTimeout    : s.ReadTimeout,
        WriteTimeout   : s.WriteTimeout,
        IdleTimeout    : s.IdleTimeout,
        MaxHeaderBytes : s.MaxHeaderBytes,
    }
}

// 设置http server参数
func (h *HttpServer)SetServerAgent(agent string) {
    h.setting.ServerAgent = agent
}

// 设置http server参数
func (h *HttpServer)SetServerRoot(root string) {
    h.setting.ServerRoot = root
}

// 绑定URI到操作函数/方法
func (h *HttpServer)BindHandle(pattern string, handler http.HandlerFunc )  {
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
func (h *HttpServer)BindHandleByMap(m map[string]http.HandlerFunc ) {
    for p, f := range m {
        h.BindHandle(p, f)
    }
}

