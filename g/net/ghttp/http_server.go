// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "sync"
    "errors"
    "strings"
    "reflect"
    "net/http"
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/net/grouter"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/gregx"
)

const (
    gHTTP_METHODS             = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
    gDEFAULT_SERVER           = "default"
    gDEFAULT_DOMAIN           = "default"
    gDEFAULT_METHOD           = "all"

    gDEFAULT_COOKIE_PATH      = "/"           // 默认path
    gDEFAULT_COOKIE_MAX_AGE   = 86400*365     // 默认cookie有效期(一年)
    gDEFAULT_SESSION_MAX_AGE  = 600           // 默认session有效期(600秒)
    gDEFAULT_SESSION_ID_NAME  = "gfsessionid" // 默认存放Cookie中的SessionId名称
)

// http server结构体
type Server struct {
    hmu           sync.RWMutex             // handlerMap互斥锁
    name          string                   // 服务名称，方便识别
    server        http.Server              // 底层http server对象
    config        ServerConfig             // 配置对象
    status        int8                     // 当前服务器状态(0：未启动，1：运行中)
    handlerMap    HandlerMap               // 所有注册的回调函数
    methodsMap    map[string]bool          // 所有支持的HTTP Method(初始化时自动填充)
    closeQueue    *gqueue.Queue            // 请求结束的关闭队列(存放的是需要异步关闭处理的*Request对象)
    hooksMap      *gmap.StringInterfaceMap // 钩子注册方法map，键值为按照注册顺序生成的glist，用于hook顺序调用
    servedCount   *gtype.Int               // 已经服务的请求数(4-8字节，不考虑溢出情况)
    cookieMaxAge  *gtype.Int               // Cookie有效期
    sessionMaxAge *gtype.Int               // Session有效期
    sessionIdName *gtype.String            // SessionId名称
    cookies       *gmap.IntInterfaceMap    // 当前服务器正在服务(请求正在执行)的Cookie(每个请求一个Cookie对象)
    sessions      *gcache.Cache            // Session内存缓存

    Router        *grouter.Router          // 路由管理对象
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
type HandlerFunc func(*Request)

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
        name          : name,
        handlerMap    : make(HandlerMap),
        methodsMap    : make(map[string]bool),
        servedCount   : gtype.NewInt(),
        closeQueue    : gqueue.New(),
        hooksMap      : gmap.NewStringInterfaceMap(),
        Router        : grouter.New(),
        cookies       : gmap.NewIntInterfaceMap(),
        sessions      : gcache.New(),
        cookieMaxAge  : gtype.NewInt(gDEFAULT_COOKIE_MAX_AGE),
        sessionMaxAge : gtype.NewInt(gDEFAULT_SESSION_MAX_AGE),
        sessionIdName : gtype.NewString(gDEFAULT_SESSION_ID_NAME),
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
    // 开启异步处理队列处理循环
    s.startCloseQueueLoop()
    // 执行端口监听
    if err := s.server.ListenAndServe(); err != nil {
        return err
    }
    s.status = 1
    return nil
}

// 生成回调方法查询的Key
func (s *Server) handlerKey(domain, method, pattern string) string {
    return strings.ToUpper(method) + ":" + pattern + "@" + strings.ToLower(domain)
}

// 设置请求处理方法
func (s *Server) setHandler(domain, method, pattern string, item HandlerItem) {
    s.hmu.Lock()
    defer s.hmu.Unlock()
    if method == gDEFAULT_METHOD {
        for v, _ := range s.methodsMap {
            s.handlerMap[s.handlerKey(domain, v, pattern)] = item
        }
    } else {
        s.handlerMap[s.handlerKey(domain, method, pattern)] = item
    }

}

// 解析pattern
func (s *Server)parsePatternForBindHandler(pattern string) (domain, method, uri string, err error) {
    uri    = pattern
    domain = gDEFAULT_DOMAIN
    method = "all"
    if array, err := gregx.MatchString(`([a-zA-Z]+):(.+)`, pattern); len(array) > 1 && err == nil {
        method  = array[1]
        pattern = array[2]
    }
    if array, err := gregx.MatchString(`(.+)@([\w\.\-]+)`, pattern); len(array) > 1 && err == nil {
        uri     = array[1]
        domain  = array[2]
    }
    if uri == "" {
        err = errors.New("invalid pattern")
    }
    return
}

// 绑定URI到操作函数/方法
// pattern的格式形如：/user/list, put:/user, delete:/user, post:/user@johng.cn
// 支持RESTful的请求格式，具体业务逻辑由绑定的处理方法来执行
func (s *Server)bindHandlerItem(pattern string, item HandlerItem) error {
    if s.status == 1 {
        return errors.New("server handlers cannot be changed while running")
    }
    domain, method, uri, err := s.parsePatternForBindHandler(pattern)
    if err != nil {
        return errors.New("invalid pattern")
    }
    s.setHandler(domain, method, uri, item)
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
        m[key] = HandlerItem{nil, "", v.Method(i).Interface().(func(*Request))}
    }
    return s.bindHandlerByMap(m)
}

// 绑定对象到URI请求处理中，会自动识别方法名称，并附加到对应的URI地址后面
// 第三个参数methods支持多个方法注册，多个方法以英文“,”号分隔，区分大小写
func (s *Server)BindObjectMethod(pattern string, obj interface{}, methods string) error {
    m := make(HandlerMap)
    for _, v := range strings.Split(methods, ",") {
        method := strings.TrimSpace(v)
        fval   := reflect.ValueOf(obj).MethodByName(method)
        if !fval.IsValid() {
            return errors.New("invalid method name:" + method)
        }
        key   := s.appendMethodNameToUriWithPattern(pattern, method)
        m[key] = HandlerItem{nil, "", fval.Interface().(func(*Request))}
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
        m[key] = HandlerItem{nil, "", v.Method(i).Interface().(func(*Request))}
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

// 这种方式绑定的控制器每一次请求都会初始化一个新的控制器对象进行处理，对应不同的请求会话
// 第三个参数methods支持多个方法注册，多个方法以英文“,”号分隔，不区分大小写
func (s *Server)BindControllerMethod(pattern string, c Controller, methods string) error {
    m    := make(HandlerMap)
    cval := reflect.ValueOf(c)
    for _, v := range strings.Split(methods, ",") {
        ctype  := reflect.ValueOf(c).Elem().Type()
        method := strings.TrimSpace(v)
        if !cval.MethodByName(method).IsValid() {
            return errors.New("invalid method name:" + method)
        }
        key    := s.appendMethodNameToUriWithPattern(pattern, method)
        m[key]  = HandlerItem{ctype, method, nil}
    }
    return s.bindHandlerByMap(m)
}

// 绑定指定的hook回调函数, hook参数的值由ghttp server设定，参数不区分大小写
// 目前hook支持：Init/Shut
func (s *Server)BindHookHandler(pattern string, hook string, handler HandlerFunc) error {
    domain, method, uri, err := s.parsePatternForBindHandler(pattern)
    if err != nil {
        return errors.New("invalid pattern")
    }
    var l *glist.List
    if method == gDEFAULT_METHOD {
        for v, _ := range s.methodsMap {
            if v := s.hooksMap.GetWithDefault(s.handlerHookKey(domain, v, uri, hook), glist.New()); v != nil {
                l = v.(*glist.List)
            }
            l.PushBack(handler)
        }
    } else {
        if v := s.hooksMap.GetWithDefault(s.handlerHookKey(domain, method, uri, hook), glist.New()); v == nil {
            l = v.(*glist.List)
        }
        l.PushBack(handler)
    }
    return nil
}

// 通过map批量绑定回调函数
func (s *Server)BindHookHandlerByMap(pattern string, hookmap map[string]HandlerFunc) error {
    for k, v := range hookmap {
        if err := s.BindHookHandler(pattern, k, v); err != nil {
            return err
        }
    }
    return nil
}

//// 绑定URI服务注册的Init回调函数，回调时按照注册顺序执行
//// Init回调调用时机为请求进入控制器之前，初始化Request对象之后
//func (s *Server)BindHookHandlerInit(pattern string, handler HandlerFunc) error {
//    return s.BindHookHandler(pattern, "Init", handler)
//}
//
//// 绑定URI服务注册的Shut回调函数，回调时按照注册顺序执行
//// Shut回调调用时机为请求执行完成之后，所有的请求资源释放之前
//func (s *Server)BindHookHandlerShut(pattern string, handler HandlerFunc) error {
//    return s.BindHookHandler(pattern, "Shut", handler)
//}

// 构造用于hooksMap检索的键名
func (s *Server)handlerHookKey(domain, method, uri, hook string) string {
    return strings.ToUpper(hook) + "^" + s.handlerKey(domain, method, uri)
}

// 获取指定hook的回调函数列表，按照注册顺序排序
func (s *Server)getHookList(domain, method, uri, hook string) []HandlerFunc {
    if v := s.hooksMap.Get(s.handlerHookKey(domain, method, uri, hook)); v != nil {
        items := v.(*glist.List).FrontAll()
        funcs := make([]HandlerFunc, len(items))
        for k, v := range items {
            funcs[k] = v.(HandlerFunc)
        }
        return funcs
    }
    return nil
}