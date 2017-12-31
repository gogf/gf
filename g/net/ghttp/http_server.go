// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package ghttp

import (
    "strings"
    "time"
    "log"
    "sync"
    "errors"
    "reflect"
    "strconv"
    "net/http"
    "crypto/tls"
    "path/filepath"
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/net/grouter"
    "sync/atomic"
)

const (
    gDEFAULT_SERVER  = "default"
    gDEFAULT_DOMAIN  = "default"
    gDEFAULT_METHOD  = "all"
    gHTTP_METHODS    = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
)

// http server结构体
type Server struct {
    hmu        sync.RWMutex    // handlerMap互斥锁
    name       string          // 服务名称，方便识别
    server     http.Server     // 底层http server对象
    config     ServerConfig    // 配置对象
    status     int8            // 当前服务器状态(0：未启动，1：运行中)
    served     uint64          // 已服务的请求数(递增)
    handlerMap HandlerMap      // 所有注册的回调函数
    methodsMap map[string]bool // 所有支持的HTTP Method
    Router     *grouter.Router // 路由管理对象
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
func GetServer(names...string) (*Server) {
    name := gDEFAULT_SERVER
    if len(names) > 0 {
        name = names[0]
    }
    if s := serverMapping.Get(name); s != nil {
        return s.(*Server)
    }
    s := &Server{
        name       : name,
        handlerMap : make(HandlerMap),
        methodsMap : make(map[string]bool),
        Router     : grouter.New(),
    }
    for _, v := range strings.Split(gHTTP_METHODS, ",") {
        s.methodsMap[v] = true
    }
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

// 服务请求数原子递增
func (s *Server) increServed() uint64 {
    return atomic.AddUint64(&s.served, 1)
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
        for v, _ := range s.methodsMap {
            s.handlerMap[s.handlerKey(domain, v, pattern)] = hitem
        }
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
        uri    = result[1]
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

// 将方法名称按照设定的规则转换为URI并附加到指定的URI后面
func (s *Server)appendMethodNameToUriWithPattern(pattern string, name string) string {
    // 检测域名后缀
    array := strings.Split(pattern, "@")
    // 分离URI(其实可能包含HTTP Method)
    uri := array[0]
    uri = strings.TrimRight(uri, "/") + "/"
    // 方法名中间存在大写字母，转换为小写URI地址以“-”号链接每个单词
    for i := 0; i < len(name); i++ {
        if i > 0 && gutil.IsLetterUpper(name[i]) {
            uri += "-"
        }
        uri += strings.ToLower(string(name[i]))
    }
    // 加上指定域名后缀
    if len(array) > 1 {
        uri += "@" + array[1]
    }
    return uri
}

// 注意该方法是直接绑定函数的内存地址，执行的时候直接执行该方法，不会存在初始化新的控制器逻辑
func (s *Server)BindHandler(pattern string, handler HandlerFunc) error {
    return s.bindHandlerItem(pattern, HandlerItem{nil, "", handler})
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (s *Server)BindObject(pattern string, obj interface{}) error {
    m := make(HandlerMap)
    v := reflect.ValueOf(obj)
    t := v.Type()
    for i := 0; i < v.NumMethod(); i++ {
        name  := t.Method(i).Name
        key   := s.appendMethodNameToUriWithPattern(pattern, name)
        m[key] = HandlerItem{nil, "", v.Method(i).Interface().(func(*Server, *ClientRequest, *ServerResponse))}
    }
    return s.bindHandlerByMap(m)
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 需要注意对象方法的定义必须按照ghttp.HandlerFunc来定义
func (s *Server)BindObjectRest(pattern string, obj interface{}) error {
    m := make(HandlerMap)
    v := reflect.ValueOf(obj)
    t := v.Type()
    for i := 0; i < v.NumMethod(); i++ {
        name := t.Method(i).Name
        if _, ok := s.methodsMap[strings.ToUpper(name)]; !ok {
            continue
        }
        key   := name + ":" + pattern
        m[key] = HandlerItem{nil, "", v.Method(i).Interface().(func(*Server, *ClientRequest, *ServerResponse))}
    }
    return s.bindHandlerByMap(m)
}

// 绑定控制器，控制器需要实现gmvc.Controller接口
// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
func (s *Server)BindController(pattern string, c Controller) error {
    // 遍历控制器，获取方法列表，并构造成uri
    m := make(HandlerMap)
    v := reflect.ValueOf(c)
    t := v.Type()
    for i := 0; i < v.NumMethod(); i++ {
        name := t.Method(i).Name
        if name == "Init" || name == "Shut" {
            continue
        }
        key   := s.appendMethodNameToUriWithPattern(pattern, name)
        m[key] = HandlerItem{v.Elem().Type(), name, nil}
    }
    return s.bindHandlerByMap(m)
}

// 绑定控制器(RESTFul)，控制器需要实现gmvc.Controller接口
// 方法会识别HTTP方法，并做REST绑定处理，例如：Post方法会绑定到HTTP POST的方法请求处理，Delete方法会绑定到HTTP DELETE的方法请求处理
// 因此只会绑定HTTP Method对应的方法，其他方法不会自动注册绑定
// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
func (s *Server)BindControllerRest(pattern string, c Controller) error {
    // 遍历控制器，获取方法列表，并构造成uri
    m := make(HandlerMap)
    v := reflect.ValueOf(c)
    t := v.Type()
    methods := make(map[string]bool)
    for _, v := range strings.Split(gHTTP_METHODS, ",") {
        methods[v] = true
    }
    for i := 0; i < v.NumMethod(); i++ {
        name := t.Method(i).Name
        if name == "Init" || name == "Shut" {
            continue
        }
        if _, ok := s.methodsMap[strings.ToUpper(name)]; !ok {
            continue
        }
        key   := name + ":" + pattern
        m[key] = HandlerItem{v.Elem().Type(), name, nil}
    }
    return s.bindHandlerByMap(m)
}

// 绑定控制器方法，pattern支持http method
// pattern的格式形如：/user/list, put:/user, delete:/user
// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
func (s *Server)BindControllerMethod(pattern string, c Controller, method string) error {
    return s.bindHandlerItem(pattern, HandlerItem{reflect.ValueOf(c).Elem().Type(), method, nil})
}
