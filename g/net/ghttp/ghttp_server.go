// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "os"
    "sync"
    "errors"
    "strings"
    "reflect"
    "runtime"
    "net/http"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/container/gqueue"
    "gitee.com/johng/gf/g/os/gspath"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/genv"
    "github.com/gorilla/websocket"
    "gitee.com/johng/gf/g/os/gtime"
    "time"
    "container/list"
)

const (
    SERVER_STATUS_STOPPED      = 0               // Server状态：停止
    SERVER_STATUS_RUNNING      = 1               // Server状态：运行
)
const (
    gHTTP_METHODS              = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
    gDEFAULT_SERVER            = "default"
    gDEFAULT_DOMAIN            = "default"
    gDEFAULT_METHOD            = "ALL"
    gDEFAULT_COOKIE_PATH       = "/"             // 默认path
    gDEFAULT_COOKIE_MAX_AGE    = 86400*365       // 默认cookie有效期(一年)
    gDEFAULT_SESSION_MAX_AGE   = 600             // 默认session有效期(600秒)
    gDEFAULT_SESSION_ID_NAME   = "gfsessionid"   // 默认存放Cookie中的SessionId名称
)

// ghttp.Server结构体
type Server struct {
    // 基本属性变量
    name             string                   // 服务名称，方便识别
    paths            *gspath.SPath            // 静态文件检索对象
    config           ServerConfig             // 配置对象
    servers          []*gracefulServer        // 底层http.Server列表
    methodsMap       map[string]bool          // 所有支持的HTTP Method(初始化时自动填充)
    servedCount      *gtype.Int               // 已经服务的请求数(4-8字节，不考虑溢出情况)，同时作为请求ID
    closeQueue       *gqueue.Queue            // 请求结束的关闭队列(存放的是需要异步关闭处理的*Request对象)
    // 服务注册相关
    handlerTree      map[string]interface{}   // 所有注册的回调函数(路由表，树型结构，哈希表+链表优先级匹配)
    handlerCache     *gcache.Cache            // 服务注册路由内存缓存
    hooksCache       *gcache.Cache            // 事件回调路由内存缓存
    // 自定义状态码回调
    hsmu             sync.RWMutex             // status handler互斥锁
    statusHandlerMap map[string]HandlerFunc   // 不同状态码下的注册处理方法(例如404状态时的处理方法)
    // COOKIE
    cookieMaxAge     *gtype.Int               // Cookie有效期
    cookies          *gmap.IntInterfaceMap    // 当前服务器正在服务(请求正在执行)的Cookie(每个请求一个Cookie对象)
    // SESSION
    sessionMaxAge    *gtype.Int               // Session有效期
    sessionIdName    *gtype.String            // SessionId名称
    sessions         *gcache.Cache            // Session内存缓存
    // 日志相关属性
    logPath          *gtype.String            // 存放日志的目录路径
    logHandler       *gtype.Interface         // 自定义日志处理回调方法
    errorLogEnabled  *gtype.Bool              // 是否开启error log
    accessLogEnabled *gtype.Bool              // 是否开启access log
    accessLogger     *glog.Logger             // access log日志对象
    errorLogger      *glog.Logger             // error log日志对象
}

// 路由对象
type Router struct {
    Uri      string       // 注册时的pattern - uri
    Method   string       // 注册时的pattern - method
    Domain   string       // 注册时的pattern - domain
    RegRule  string       // 路由规则解析后对应的正则表达式
    RegNames []string     // 路由规则解析后对应的变量名称数组
    Priority int          // 优先级，用于链表排序，值越大优先级越高
}

// 域名、URI与回调函数的绑定记录表
type handlerMap  map[string]*handlerItem

// http回调函数注册信息
type handlerItem struct {
    ctype    reflect.Type // 控制器类型(反射类型)
    fname    string       // 回调方法名称
    faddr    HandlerFunc  // 准确的执行方法内存地址(与以上两个参数二选一)
}

// 路由注册项(这里使用了非并发安全的list.List，因为该对象的使用统一是由htmu互斥锁保证并发安全)
type handlerRegisterItem struct {
    handler    *handlerItem             // 准确的执行方法内存地址
    hooks      map[string]*list.List   // 当前的事件回调注册，键名为事件名称，键值为事件执行方法链表
    router     *Router                  // 注册时绑定的路由对象
}

// 根据特定URL.Path解析后的路由检索结果项
type handlerParsedItem struct {
    item    *handlerRegisterItem        // 路由注册项
    values   map[string][]string        // 特定URL.Path的Router解析参数
}

// HTTP注册函数
type HandlerFunc func(r *Request)

// 文件描述符map
type listenerFdMap map[string]string


// Server表，用以存储和检索名称与Server对象之间的关联关系
var serverMapping    = gmap.NewStringInterfaceMap()

// 正常运行的Server数量，如果没有运行、失败或者全部退出，那么该值为0
var serverRunning    = gtype.NewInt()

// Web Socket默认配置
var wsUpgrader       = websocket.Upgrader{}

// Web Server已完成服务事件通道，当有事件时表示服务完成，当前进程退出
var doneChan         = make(chan struct{}, 1000)

// 用于服务进程初始化，只能初始化一次，采用“懒初始化”(在server运行时才初始化)
var serverProcInited = gtype.NewBool()

// Web Server进程初始化.
// 注意该方法不能放置于包初始化方法init中，不使用ghttp.Server的功能便不能初始化对应的协程goroutine逻辑.
func serverProcInit() {
    if serverProcInited.Val() {
        return
    }
    serverProcInited.Set(true)
    // 如果是完整重启，那么需要等待主进程销毁后，才开始执行监听，防止端口冲突
    if genv.Get(gADMIN_ACTION_RESTART_ENVKEY) != "" {
        if p, e := os.FindProcess(gproc.PPid()); e == nil {
            p.Kill()
            p.Wait()
        } else {
            glog.Error(e)
        }
    }

    // 信号量管理操作监听
    go handleProcessSignal()
    // 异步监听进程间消息
    go handleProcessMessage()
}

// 获取/创建一个默认配置的HTTP Server(默认监听端口是80)
// 单例模式，请保证name的唯一性
func GetServer(name...interface{}) (*Server) {
    sname := gDEFAULT_SERVER
    if len(name) > 0 {
        sname = gconv.String(name[0])
    }
    if s := serverMapping.Get(sname); s != nil {
        return s.(*Server)
    }
    s := &Server {
        name             : sname,
        paths            : gspath.New(),
        servers          : make([]*gracefulServer, 0),
        methodsMap       : make(map[string]bool),
        statusHandlerMap : make(map[string]HandlerFunc),
        handlerTree      : make(map[string]interface{}),
        handlerCache     : gcache.New(),
        hooksCache       : gcache.New(),
        cookies          : gmap.NewIntInterfaceMap(),
        sessions         : gcache.New(),
        cookieMaxAge     : gtype.NewInt(gDEFAULT_COOKIE_MAX_AGE),
        sessionMaxAge    : gtype.NewInt(gDEFAULT_SESSION_MAX_AGE),
        sessionIdName    : gtype.NewString(gDEFAULT_SESSION_ID_NAME),
        servedCount      : gtype.NewInt(),
        closeQueue       : gqueue.New(),
        logPath          : gtype.NewString(),
        accessLogEnabled : gtype.NewBool(),
        errorLogEnabled  : gtype.NewBool(true),
        accessLogger     : glog.New(),
        errorLogger      : glog.New(),
        logHandler       : gtype.NewInterface(),
    }
    s.errorLogger.SetBacktraceSkip(4)
    s.accessLogger.SetBacktraceSkip(4)
    // 设置路由解析缓存上限，使用LRU进行缓存淘汰
    s.hooksCache.SetCap(10000)
    s.handlerCache.SetCap(10000)
    for _, v := range strings.Split(gHTTP_METHODS, ",") {
        s.methodsMap[v] = true
    }
    s.SetConfig(defaultServerConfig)
    serverMapping.Set(sname, s)
    return s
}

// 作为守护协程异步执行(当同一进程中存在多个Web Server时，需要采用这种方式执行)
// 需要结合Wait方式一起使用
func (s *Server) Start() error {
    serverProcInit()

    // 如果设置了静态文件目录，那么严格按照静态文件目录进行检索
    // 否则，默认使用当前可执行文件目录，并且如果是开发环境，默认也会添加main包的源码目录路径做为二级检索
    if s.config.ServerRoot != "" {
        s.paths.Set(s.config.ServerRoot)
    } else {
        s.paths.Set(gfile.SelfDir())
        if p := gfile.MainPkgPath(); gfile.Exists(p) {
            s.paths.Add(p)
        }
    }

    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server is already running")
    }
    // 底层http server配置
    if s.config.Handler == nil {
        s.config.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }

    // 启动http server
    reloaded := false
    fdMapStr := genv.Get(gADMIN_ACTION_RELOAD_ENVKEY)
    if len(fdMapStr) > 0 {
        sfm := bufferToServerFdMap([]byte(fdMapStr))
        if v, ok := sfm[s.name]; ok {
            s.startServer(v)
            reloaded = true
        }
    }
    if !reloaded {
        s.startServer(nil)
    }

    // 如果是子进程，那么服务开启后通知父进程销毁
    if gproc.IsChild() {
        gtime.SetTimeout(2*time.Second, func() {
            gproc.Send(gproc.PPid(), []byte("exit"), gADMIN_GPROC_COMM_GROUP)
        })
    }

    // 开启异步关闭队列处理循环
    s.startCloseQueueLoop()
    return nil
}

// 阻塞执行监听
func (s *Server) Run() error {
    if err := s.Start(); err != nil {
        return err
    }
    // 阻塞等待服务执行完成
    <- doneChan

    glog.Printfln("%d: all servers shutdown", gproc.Pid())
    return nil
}


// 阻塞等待所有Web Server停止，常用于多Web Server场景，以及需要将Web Server异步运行的场景
// 这是一个与进程相关的方法
func Wait() {
    // 阻塞等待服务执行完成
    <- doneChan

    glog.Printfln("%d: all servers shutdown", gproc.Pid())
}


// 开启底层Web Server执行
func (s *Server) startServer(fdMap listenerFdMap) {
    var httpsEnabled bool
    if len(s.config.HTTPSCertPath) > 0 && len(s.config.HTTPSKeyPath) > 0 {
        // ================
        // HTTPS
        // ================
        if len(s.config.HTTPSAddr) == 0 {
            if len(s.config.Addr) > 0 {
                s.config.HTTPSAddr = s.config.Addr
                s.config.Addr      = ""
            } else {
                s.config.HTTPSAddr = gDEFAULT_HTTPS_ADDR
            }
        }
        httpsEnabled = len(s.config.HTTPSAddr) > 0
        var array []string
        if v, ok := fdMap["https"]; ok && len(v) > 0 {
            array = strings.Split(v, ",")
        } else {
            array = strings.Split(s.config.HTTPSAddr, ",")
        }
        for _, v := range array {
            if len(v) == 0 {
                continue
            }
            fd    := 0
            addr  := v
            array := strings.Split(v, "#")
            if len(array) > 1 {
                addr = array[0]
                // windows系统不支持文件描述符传递socket通信平滑交接，因此只能完整重启
                if runtime.GOOS != "windows" {
                    fd = gconv.Int(array[1])
                }
            }
            if fd > 0 {
                s.servers = append(s.servers, s.newGracefulServer(addr, fd))
            } else {
                s.servers = append(s.servers, s.newGracefulServer(addr))
            }
            s.servers[len(s.servers) - 1].isHttps = true
        }
    }
    // ================
    // HTTP
    // ================
    // 当HTTPS服务未启用时，默认HTTP地址才会生效
    if !httpsEnabled && len(s.config.Addr) == 0 {
        s.config.Addr = gDEFAULT_HTTP_ADDR
    }
    var array []string
    if v, ok := fdMap["http"]; ok && len(v) > 0 {
        array = strings.Split(v, ",")
    } else {
        array = strings.Split(s.config.Addr, ",")
    }
    for _, v := range array {
        if len(v) == 0 {
            continue
        }
        fd    := 0
        addr  := v
        array := strings.Split(v, "#")
        if len(array) > 1 {
            addr = array[0]
            // windows系统不支持文件描述符传递socket通信平滑交接，因此只能完整重启
            if runtime.GOOS != "windows" {
                fd = gconv.Int(array[1])
            }
        }
        if fd > 0 {
            s.servers = append(s.servers, s.newGracefulServer(addr, fd))
        } else {
            s.servers = append(s.servers, s.newGracefulServer(addr))
        }
    }
    // 开始执行异步监听
    for _, v := range s.servers {
        go func(server *gracefulServer) {
            serverRunning.Add(1)
            err := (error)(nil)
            if server.isHttps {
                err = server.ListenAndServeTLS(s.config.HTTPSCertPath, s.config.HTTPSKeyPath)
            } else {
                err = server.ListenAndServe()
            }
            serverRunning.Add(-1)
            // 如果非关闭错误，那么提示报错，否则认为是正常的服务关闭操作
            if err != nil && !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
                glog.Error(err)
            }
            // 如果所有异步的Server都已经停止，那么主Server就可以退出了
            if serverRunning.Val() < 1 {
                doneChan <- struct{}{}
            }
        }(v)
    }
}

// 获取当前服务器的状态
func (s *Server) Status() int {
    if serverRunning.Val() > 0 {
        return SERVER_STATUS_RUNNING
    }
    return SERVER_STATUS_STOPPED
}

// 获取当前监听的文件描述符信息，构造成map返回
func (s *Server) getListenerFdMap() map[string]string {
    m := map[string]string {
        "https" : "",
        "http"  : "",
    }
    // s.servers是从HTTPS到HTTP优先级遍历，解析的时候也应当按照这个顺序读取fd
    for _, v := range s.servers {
        str := v.addr + "#" + gconv.String(v.Fd()) + ","
        if v.isHttps {
            m["https"] += str
        } else {
            m["http"]  += str
        }
    }
    // 去掉末尾的","号
    if len(m["https"]) > 0 {
        m["https"] = m["https"][0 : len(m["https"]) - 1]
    }
    if len(m["http"]) > 0 {
        m["http"] = m["http"][0 : len(m["http"]) - 1]
    }

    return m
}
