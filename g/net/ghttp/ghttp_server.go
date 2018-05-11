// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "os"
    "net"
    "sync"
    "errors"
    "strings"
    "reflect"
    "net/http"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/container/gqueue"
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
    config           ServerConfig             // 配置对象
    status           int8                     // 当前服务器状态(0：未启动，1：运行中)
    servers          []*gracefulServer        // 底层http.Server列表
    methodsMap       map[string]bool          // 所有支持的HTTP Method(初始化时自动填充)
    servedCount      *gtype.Int               // 已经服务的请求数(4-8字节，不考虑溢出情况)，同时作为请求ID
    closeQueue       *gqueue.Queue            // 请求结束的关闭队列(存放的是需要异步关闭处理的*Request对象)
    signalQueue      chan os.Signal           // 终端命令行监听队列
    // 服务注册相关
    hmmu             sync.RWMutex             // handler互斥锁
    hmcmu            sync.RWMutex             // handlerCache互斥锁
    handlerMap       HandlerMap               // 所有注册的回调函数(静态匹配)
    handlerTree      map[string]interface{}   // 所有注册的回调函数(动态匹配，树型+链表优先级匹配)
    handlerCache     *gcache.Cache            // 服务注册路由内存缓存
    // 事件回调注册
    hhmu             sync.RWMutex             // hooks互斥锁
    hhcmu            sync.RWMutex             // hooksCache互斥锁
    hooksTree        map[string]interface{}   // 所有注册的事件回调函数(动态匹配，树型+链表优先级匹配)
    hooksCache       *gcache.Cache            // 回调事件注册路由内存缓存
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

// 域名、URI与回调函数的绑定记录表
type HandlerMap  map[string]*HandlerItem

// 路由对象
type Router struct {
    Uri      string       // 注册时的pattern - uri
    Method   string       // 注册时的pattern - method
    Domain   string       // 注册时的pattern - domain
    Priority int          // 优先级，用于链表排序，值越大优先级越高
}

// http回调函数注册信息
type HandlerItem struct {
    ctype    reflect.Type // 控制器类型
    fname    string       // 回调方法名称
    faddr    HandlerFunc  // 准确的执行方法内存地址(与以上两个参数二选一)
    router   *Router      // 注册时绑定的路由对象
}

// http注册函数
type HandlerFunc func(r *Request)

// 文件描述符map
type listenerFdMap map[string]string

// Server表，用以存储和检索名称与Server对象之间的关联关系
var serverMapping = gmap.NewStringInterfaceMap()

// Web Server多进程管理器
var procManager   = gproc.NewManager()

// Web Server开始执行事件通道，由于同一个进程支持多Server，因此该通道为非阻塞
var runChan       = make(chan struct{}, 100000)
// 已完成的run消息队列
var doneChan      = make(chan struct{}, 100000)

// Web Server进程初始化
func init() {
    go func() {
        // 等待run消息(Run方法调用)
        <- runChan
        // 主进程只负责创建子进程
        if !gproc.IsChild() {
            sendProcessMsg(os.Getpid(), gMSG_START, nil)
        }
        // 开启进程消息监听处理
        handleProcessMsg()
        // 服务执行完成，需要退出
        doneChan <- struct{}{}
    }()
}

// 获取所有Web Server的文件描述符map
func getServerFdMap() map[string]listenerFdMap {
    sfm := make(map[string]listenerFdMap)
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for k, v := range m {
            sfm[k] = v.(*Server).getListenerFdMap()
        }
    })
    return sfm
}

// 二进制转换为FdMap
func bufferToServerFdMap(buffer []byte) map[string]listenerFdMap {
    sfm := make(map[string]listenerFdMap)
    if len(buffer) > 0 {
        j, _ := gjson.LoadContent(buffer, "json")
        for k, _ := range j.ToMap() {
            m := make(map[string]string)
            for k, v := range j.GetMap(k) {
                m[k] = gconv.String(v)
            }
            sfm[k] = m
        }
    }
    return sfm
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
        servers          : make([]*gracefulServer, 0),
        methodsMap       : make(map[string]bool),
        handlerMap       : make(HandlerMap),
        statusHandlerMap : make(map[string]HandlerFunc),
        handlerTree      : make(map[string]interface{}),
        hooksTree        : make(map[string]interface{}),
        handlerCache     : gcache.New(),
        hooksCache       : gcache.New(),
        cookies          : gmap.NewIntInterfaceMap(),
        sessions         : gcache.New(),
        cookieMaxAge     : gtype.NewInt(gDEFAULT_COOKIE_MAX_AGE),
        sessionMaxAge    : gtype.NewInt(gDEFAULT_SESSION_MAX_AGE),
        sessionIdName    : gtype.NewString(gDEFAULT_SESSION_ID_NAME),
        servedCount      : gtype.NewInt(),
        closeQueue       : gqueue.New(),
        signalQueue      : make(chan os.Signal),
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

// 阻塞执行监听
func (s *Server) Run() error {
    runChan <- struct{}{}

    if !gproc.IsChild() {
        <- doneChan
        return nil
    }

    if s.status == 1 {
        return errors.New("server is already running")
    }

    // 底层http server配置
    if s.config.Handler == nil {
        s.config.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }

    // 开启异步关闭队列处理循环
    s.startCloseQueueLoop()

    // 阻塞等待服务执行完成
    <- doneChan

    glog.Printfln("web server pid %d exit successfully", gproc.Pid())
    return nil
}

// 开启底层Web Server执行
func (s *Server) startServer(fdMap listenerFdMap) {
    // 开始执行底层Web Server创建，端口监听
    var server *gracefulServer
    if len(s.config.HTTPSCertPath) > 0 && len(s.config.HTTPSKeyPath) > 0 {
        // HTTPS
        if len(s.config.HTTPSAddr) == 0 {
            if len(s.config.Addr) > 0 {
                s.config.HTTPSAddr = s.config.Addr
            } else {
                s.config.HTTPSAddr = gDEFAULT_HTTPS_ADDR
            }
        }
        var array []string
        var isFd  bool
        if v, ok := fdMap["https"]; ok && len(v) > 0 {
            isFd  = true
            array = strings.Split(v, ",")
        } else {
            array = strings.Split(s.config.HTTPSAddr, ",")
        }

        for _, v := range array {
            go func(item string) {
                if isFd {
                    tArray := strings.Split(item, "#")
                    server  = s.newGracefulServer(tArray[0], gconv.Int(tArray[1]))
                } else {
                    server  = s.newGracefulServer(item)
                }
                s.servers = append(s.servers, server)
                glog.Printfln("https server started listening on %s", server.addr)
                if err := server.ListenAndServeTLS(s.config.HTTPSCertPath, s.config.HTTPSKeyPath); err != nil {
                    // 如果非关闭错误，那么提示报错，否则认为是正常的服务关闭操作
                    if !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
                        glog.Error(err)
                    }
                }
            }(v)
        }
    }
    // HTTP
    if s.servedCount.Val() == 0 && len(s.config.Addr) == 0 {
        s.config.Addr = gDEFAULT_HTTP_ADDR
    }
    var array []string
    var isFd  bool
    if v, ok := fdMap["http"]; ok && len(v) > 0 {
        isFd  = true
        array = strings.Split(v, ",")
    } else {
        array = strings.Split(s.config.Addr, ",")
    }
    for _, v := range array {
        go func(item string) {
            if isFd {
                tArray := strings.Split(item, "#")
                server  = s.newGracefulServer(tArray[0], gconv.Int(tArray[1]))
            } else {
                server  = s.newGracefulServer(item)
            }
            s.servers = append(s.servers, server)
            glog.Printfln("http server started listening on %s", server.addr)
            if err := server.ListenAndServe(); err != nil {
                // 如果非关闭错误，那么提示报错，否则认为是正常的服务关闭操作
                if !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
                    glog.Error(err)
                }
            }
        }(v)
    }

    s.status = 1
}

// 重启Web Server
func (s *Server) Restart() {
    // 如果是主进程，那么向所有子进程发送重启信号
    if !gproc.IsChild() {
        procManager.Send(formatMsgBuffer(gMSG_RESTART, nil))
        return
    }
    sendProcessMsg(gproc.Pid(), gMSG_RESTART, nil)
}

// 关闭Web Server
func (s *Server) Shutdown() {
    // 如果是主进程，那么向所有子进程发送关闭信号
    if !gproc.IsChild() {
        procManager.Send(formatMsgBuffer(gMSG_SHUTDOWN, nil))
        return
    }
    for _, v := range s.servers {
        v.shutdown()
    }
}

// 获取当前监听的文件描述符信息，构造成map返回
func (s *Server) getListenerFdMap() map[string]string {
    m := map[string]string {
        "http"  : "",
        "https" : "",
    }
    for _, v := range s.servers {
        if f, e := v.listener.(*net.TCPListener).File(); e == nil {
            str := v.addr + "#" + gconv.String(f.Fd()) + ","
            if v.isHttps {
                m["https"] += str
            } else {
                m["http"]  += str
            }
        } else {
            glog.Errorfln("failed to get listener file: %v", e)
        }
    }
    if len(m["http"]) > 0 {
        m["http"] = m["http"][0 : len(m["http"]) - 1]
    }
    if len(m["https"]) > 0 {
        m["https"] = m["https"][0 : len(m["https"]) - 1]
    }
    return m
}

// 清空当前的handlerCache
func (s *Server) clearHandlerCache() {
    s.hmcmu.Lock()
    defer s.hmcmu.Unlock()
    s.handlerCache.Close()
    s.handlerCache = gcache.New()
}

// 清空当前的hooksCache
func (s *Server) clearHooksCache() {
    s.hhcmu.Lock()
    defer s.hhcmu.Unlock()
    s.hooksCache.Close()
    s.hooksCache = gcache.New()
}