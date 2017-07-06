package g

import (
    "net/http"
    "time"
    "crypto/tls"
    "log"
    //"fmt"
    "strings"
    "path/filepath"
)
/**
    @todo 静态文件的处理性能比Nginx稍弱，不能使用标准库方法，需自行处理
 */
// 全局http封装对象
var Http gHttp

// http 结构体
type gHttp struct {
    Server GHttpServer
}

// http server结构体
type GHttpServer struct {
    server     http.Server
    setting    GHttpServerSetting
    handlerMap map[string]http.HandlerFunc
}

// HTTP Server 设置结构体
type GHttpServerSetting struct {
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
var defaultHttpServerSetting = GHttpServerSetting {
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
func (h GHttpServer)New() (*GHttpServer) {
    return h.NewBySetting(defaultHttpServerSetting)
}

// 创建一个HTTP Server，返回指针
func (h GHttpServer)NewByAddr(addr string) (*GHttpServer) {
    setting     := defaultHttpServerSetting
    setting.Addr = addr
    return h.NewBySetting(setting)
}

// 创建一个HTTP Server
func (h GHttpServer)NewByAddrRoot(addr string, root string) (*GHttpServer) {
    setting           := defaultHttpServerSetting
    setting.Addr       = addr
    setting.ServerRoot = root
    return h.NewBySetting(setting)
}

// 根据输入配置创建一个http server对象
func (h GHttpServer)NewBySetting(s GHttpServerSetting) (*GHttpServer) {
    var server GHttpServer
    server.SetSetting(s)
    return &server
}

// 执行
func (h *GHttpServer)Run() error {
    err := h.server.ListenAndServe()
    if err != nil {
        panic(err)
    }
    return err
}

// 默认HTTP Server处理入口
func (h *GHttpServer)defaultHttpHandle(w http.ResponseWriter, r *http.Request) {
    if f, ok := h.handlerMap[r.URL.String()]; ok {
        f(w, r)
    } else {
        h.serveFile(w, r)
    }
}

// 处理静态文件请求
func (h *GHttpServer)serveFile(w http.ResponseWriter, r *http.Request) {
    uri := r.URL.String()
    if h.setting.ServerRoot != "" {
        path := strings.TrimRight(h.setting.ServerRoot, string(filepath.Separator))
        path  = path + uri
        // fmt.Println(path)
        if (File.Exists(path)) {
            http.ServeFile(w, r, path)
        } else {
            http.NotFound(w, r)
        }
    } else {
        panic("http server root is empty while handling static files request")
    }
}

// 获取默认的http server设置
func (h GHttpServer)GetDefaultSetting() GHttpServerSetting {
    return defaultHttpServerSetting
}

// http server setting设置
func (h *GHttpServer)SetSetting(s GHttpServerSetting) {
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
func (h *GHttpServer)SetServerAgent(agent string) {
    h.setting.ServerAgent = agent
}

// 设置http server参数
func (h *GHttpServer)SetServerRoot(root string) {
    h.setting.ServerRoot = root
}

// 绑定URI到操作函数/方法
func (h *GHttpServer)BindHandle(pattern string, handler http.HandlerFunc )  {
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
func (h *GHttpServer)BindHandleByMap(m map[string]http.HandlerFunc ) {
    for p, f := range m {
        h.BindHandle(p, f)
    }
}

