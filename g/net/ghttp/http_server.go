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
    "reflect"
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/container/gmap"
    "strconv"
)

const (
    gDEFAULT_DOMAIN = "default"
    gDEFAULT_METHOD = "all"
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

// 域名、URI与回调函数的绑定记录表
type HandlerMap  map[string]HandlerItem

// http回调函数注册信息
type HandlerItem struct {
    ctype reflect.Type // 控制器类型
    fname string       // 回调方法名称
    faddr HandlerFunc  // 准确的执行方法内存地址(与以上两个参数二选一)
}

// http注册函数
type HandlerFunc func(*Server, *ClientRequest, *ServerResponse)

// Server表，用以存储和检索名称与Server对象之间的关联关系
var serverMapping = gmap.NewStringInterfaceMap()

// 获取/创建一个默认配置的HTTP Server(默认监听端口是80)
// 单例模式，请保证name的唯一性
func GetServer(name string) (*Server) {
    if s := serverMapping.Get(name); s != nil {
        return s.(*Server)
    }
    s           := &Server{}
    s.name       = name
    s.handlerMap = make(HandlerMap)
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

// 获取
func (s *Server) GetName() string {
    return s.name
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

// 设置http server参数 - Port
func (s *Server)SetPort(port int) error {
    if s.status == 1 {
        return errors.New("server config cannot be changed while running")
    }
    s.config.Addr = ":" + strconv.Itoa(port)
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

// 生成回调方法查询的Key
func (s *Server) handlerKey(domain, method, pattern string) string {
    return strings.ToUpper(method) + ":" + pattern + "@" + strings.ToLower(domain)
}

// 设置请求处理方法
func (s *Server) setHandler(domain, method, pattern string, hitem HandlerItem) {
    s.hmu.Lock()
    defer s.hmu.Unlock()
    if method == gDEFAULT_METHOD {
        s.handlerMap[s.handlerKey(domain, "GET",     pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "PUT",     pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "POST",    pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "DELETE",  pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "PATCH",   pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "HEAD",    pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "CONNECT", pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "OPTIONS", pattern)] = hitem
        s.handlerMap[s.handlerKey(domain, "TRACE",   pattern)] = hitem
    } else {
        s.handlerMap[s.handlerKey(domain, method, pattern)] = hitem
    }
}

// 查询请求处理方法
func (s *Server) getHandler(domain, method, pattern string) *HandlerItem {
    s.hmu.RLock()
    defer s.hmu.RUnlock()
    key := s.handlerKey(domain, method, pattern)
    if f, ok := s.handlerMap[key]; ok {
        return &f
    }
    return nil
}

// 绑定URI到操作函数/方法
// pattern的格式形如：/user/list, put:/user, delete:/user, post:/user@johng.cn
// 支持RESTful的请求格式，具体业务逻辑由绑定的处理方法来执行
func (s *Server)bindHandlerItem(pattern string, hitem HandlerItem) error {
    if s.status == 1 {
        return errors.New("server handlers cannot be changed while running")
    }
    uri    := ""
    domain := gDEFAULT_DOMAIN
    method := "all"
    result := strings.Split(pattern, "@")
    if len(result) > 1 {
        domain = result[1]
    }
    result  = strings.Split(result[0], ":")
    if len(result) > 1 {
        method = result[0]
        uri    = result[0]
    } else {
        uri    = result[0]
    }
    if uri == "" {
        return errors.New("invalid pattern")
    }
    s.setHandler(domain, method, uri, hitem)
    return nil
}

// 通过映射数组绑定URI到操作函数/方法
func (s *Server)bindHandlerByMap(m HandlerMap) error {
    for p, h := range m {
        if err := s.bindHandlerItem(p, h); err != nil {
            return err
        }
    }
    return nil
}

// 注意该方法是直接绑定方法的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (s *Server)BindHandler(pattern string, handler HandlerFunc) error {
    return s.bindHandlerItem(pattern, HandlerItem{nil, "", handler})
}

// 绑定控制器，控制器需要实现gmvc.Controller接口
// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
func (s *Server)BindController(uri string, c Controller) error {
    // 遍历控制器，获取方法列表，并构造成uri
    m := make(HandlerMap)
    v := reflect.ValueOf(c)
    t := v.Type()
    for i := 0; i < v.NumMethod(); i++ {
        key  := strings.TrimRight(uri, "/") + "/"
        name := t.Method(i).Name
        if name == "Init" || name == "Shut" {
            continue
        }
        for i := 0; i < len(name); i++ {
            if i > 0 && gutil.IsLetterUpper(name[i]) {
                key += "-"
            }
            key += strings.ToLower(string(name[i]))
        }
        m[key] = HandlerItem{v.Elem().Type(), name, nil}
    }
    return s.bindHandlerByMap(m)
}

// 绑定方法，pattern支持http method
// pattern的格式形如：/user/list, put:/user, delete:/user
// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
func (s *Server)BindControllerMethod(pattern string, c Controller, method string) error {
    return s.bindHandlerItem(pattern, HandlerItem{reflect.ValueOf(c).Elem().Type(), method, nil})
}
