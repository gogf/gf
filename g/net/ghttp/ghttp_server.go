// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "bytes"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/os/gtimer"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/third/github.com/gorilla/websocket"
    "gitee.com/johng/gf/third/github.com/olekukonko/tablewriter"
    "net/http"
    "os"
    "reflect"
    "runtime"
    "strings"
    "sync"
    "time"
)

type (
    // Server结构体
    Server struct {
        // 基本属性变量
        name             string                           // 服务名称，方便识别
        config           ServerConfig                     // 配置对象
        servers          []*gracefulServer                // 底层http.Server列表
        methodsMap       map[string]struct{}              // 所有支持的HTTP Method(初始化时自动填充)
        servedCount      *gtype.Int                       // 已经服务的请求数(4-8字节，不考虑溢出情况)，同时作为请求ID
        // 服务注册相关
        serveTree        map[string]interface{}           // 所有注册的服务回调函数(路由表，树型结构，哈希表+链表优先级匹配)
        hooksTree        map[string]interface{}           // 所有注册的事件回调函数(路由表，树型结构，哈希表+链表优先级匹配)
        serveCache       *gcache.Cache                    // 服务注册路由内存缓存
        hooksCache       *gcache.Cache                    // 事件回调路由内存缓存
        routesMap        map[string][]registeredRouteItem // 已经注册的路由及对应的注册方法文件地址(用以路由重复注册判断)
        // 自定义状态码回调
        hsmu             sync.RWMutex                     // status handler互斥锁
        statusHandlerMap map[string]HandlerFunc           // 不同状态码下的注册处理方法(例如404状态时的处理方法)
        // SESSION
        sessions         *gcache.Cache                    // Session内存缓存
        // Logger
        logger           *glog.Logger                     // 日志管理对象
    }

    // 路由对象
    Router struct {
        Uri      string       // 注册时的pattern - uri
        Method   string       // 注册时的pattern - method
        Domain   string       // 注册时的pattern - domain
        RegRule  string       // 路由规则解析后对应的正则表达式
        RegNames []string     // 路由规则解析后对应的变量名称数组
        Priority int          // 优先级，用于链表排序，值越大优先级越高
    }

    // http回调函数注册信息
    handlerItem struct {
        name     string       // 注册的方法名称信息
        rtype    int          // 注册方式(执行对象/回调函数/控制器)
        ctype    reflect.Type // 控制器类型(反射类型)
        fname    string       // 回调方法名称
        faddr    HandlerFunc  // 准确的执行方法内存地址(与以上两个参数二选一)
        finit    HandlerFunc  // 初始化请求回调方法(执行对象注册方式下有效)
        fshut    HandlerFunc  // 完成请求回调方法(执行对象注册方式下有效)
        router   *Router      // 注册时绑定的路由对象
    }

    // 根据特定URL.Path解析后的路由检索结果项
    handlerParsedItem struct {
        handler  *handlerItem         // 路由注册项
        values   map[string][]string  // 特定URL.Path的Router解析参数
    }

    // 已注册的路由项
    registeredRouteItem struct {
        file     string               // 文件路径及行数地址
        handler  *handlerItem         // 路由注册项
    }

    // pattern与回调函数的绑定map
    handlerMap    = map[string]*handlerItem

    // HTTP注册函数
    HandlerFunc   = func(r *Request)

    // 文件描述符map
    listenerFdMap = map[string]string
)

const (
    SERVER_STATUS_STOPPED      = 0               // Server状态：停止
    SERVER_STATUS_RUNNING      = 1               // Server状态：运行
    HOOK_BEFORE_SERVE          = "BeforeServe"
    HOOK_AFTER_SERVE           = "AfterServe"
    HOOK_BEFORE_OUTPUT         = "BeforeOutput"
    HOOK_AFTER_OUTPUT          = "AfterOutput"
    HOOK_BEFORE_CLOSE          = "BeforeClose"
    HOOK_AFTER_CLOSE           = "AfterClose"

    gHTTP_METHODS              = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
    gDEFAULT_SERVER            = "default"
    gDEFAULT_DOMAIN            = "default"
    gDEFAULT_METHOD            = "ALL"
    gROUTE_REGISTER_HANDLER    = 1
    gROUTE_REGISTER_OBJECT     = 2
    gROUTE_REGISTER_CONTROLLER = 3
    gEXCEPTION_EXIT            = "exit"
    gEXCEPTION_EXIT_ALL        = "exit_all"
    gEXCEPTION_EXIT_HOOK       = "exit_hook"
)

var (
    // Server表，用以存储和检索名称与Server对象之间的关联关系
    serverMapping    = gmap.NewStringInterfaceMap()

    // 正常运行的Server数量，如果没有运行、失败或者全部退出，那么该值为0
    serverRunning    = gtype.NewInt()

    // Web Socket默认配置
    wsUpgrader       = websocket.Upgrader {
        // 默认允许WebSocket请求跨域，权限控制可以由业务层自己负责，灵活度更高
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    // Web Server已完成服务事件通道，当有事件时表示服务完成，当前进程退出
    doneChan         = make(chan struct{}, 1000)

    // 用于服务进程初始化，只能初始化一次，采用“懒初始化”(在server运行时才初始化)
    serverProcInited = gtype.NewBool()
)


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
        servers          : make([]*gracefulServer, 0),
        methodsMap       : make(map[string]struct{}),
        statusHandlerMap : make(map[string]HandlerFunc),
        serveTree        : make(map[string]interface{}),
        hooksTree        : make(map[string]interface{}),
        serveCache       : gcache.New(),
        hooksCache       : gcache.New(),
        routesMap        : make(map[string][]registeredRouteItem),
        sessions         : gcache.New(),
        servedCount      : gtype.NewInt(),
        logger           : glog.New(),
    }
    // 日志的标准输出默认关闭，但是错误信息会特殊处理
    s.logger.SetStdPrint(false)
    for _, v := range strings.Split(gHTTP_METHODS, ",") {
        s.methodsMap[v] = struct{}{}
    }
    // 初始化时使用默认配置
    s.SetConfig(defaultServerConfig)
    // 记录到全局ServerMap中
    serverMapping.Set(sname, s)
    return s
}

// 作为守护协程异步执行(当同一进程中存在多个Web Server时，需要采用这种方式执行)
// 需要结合Wait方式一起使用
func (s *Server) Start() error {
    // 服务进程初始化，只会初始化一次
    serverProcInit()

    // 当前Web Server状态判断
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server is already running")
    }

    // 底层http server配置
    if s.config.Handler == nil {
        s.config.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }
    // 不允许访问的路由注册(使用HOOK实现)
    if s.config.DenyRoutes != nil {
        for _, v := range s.config.DenyRoutes {
            s.BindHookHandler(v, HOOK_BEFORE_SERVE, func(r *Request) {
                r.Response.WriteStatus(403)
                r.Exit()
            })
        }
    }

    // gzip压缩文件类型
    //if s.config.GzipContentTypes != nil {
    //    for _, v := range s.config.GzipContentTypes {
    //        s.gzipMimesMap[v] = struct{}{}
    //    }
    //}

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
        gtimer.SetTimeout(2*time.Second, func() {
            if err := gproc.Send(gproc.PPid(), []byte("exit"), gADMIN_GPROC_COMM_GROUP); err != nil {
                panic(err)
            }
        })
    }
    // 是否处于开发环境
    if gfile.MainPkgPath() != "" {
        glog.Debug("GF notices that you're in develop environment, so error logs are auto enabled to stdout.")
    }

    // 打印展示路由表
    s.DumpRoutesMap()
    return nil
}

// 打印展示路由表
func (s *Server) DumpRoutesMap() {
    if s.config.DumpRouteMap && len(s.routesMap) > 0 {
        // (等待一定时间后)当所有框架初始化信息打印完毕之后才打印路由表信息
        gtimer.SetTimeout(50*time.Millisecond, func() {
            glog.Header(false).Println(fmt.Sprintf("\n%s\n", s.GetRouteMap()))
        })
    }
}

// 获得路由表(格式化字符串)
func (s *Server) GetRouteMap() string {
    type tableItem struct {
        hook     string
        domain   string
        method   string
        route    string
        handler  string
        priority int
    }

    buf   := bytes.NewBuffer(nil)
    table := tablewriter.NewWriter(buf)
    table.SetHeader([]string{"SERVER", "ADDRESS", "DOMAIN", "METHOD", "P", "ROUTE", "HANDLER", "HOOK"})
    table.SetRowLine(true)
    table.SetBorder(false)
    table.SetCenterSeparator("|")

    m := make(map[string]*garray.SortedArray)
    for k, registeredItems := range s.routesMap {
        array, _ := gregex.MatchString(`(.*?)%([A-Z]+):(.+)@(.+)`, k)
        for index, registeredItem := range registeredItems {
            item := &tableItem {
                hook     : array[1],
                domain   : array[4],
                method   : array[2],
                route    : array[3],
                handler  : registeredItem.handler.name,
                priority : len(registeredItems) - index - 1,
            }
            if _, ok := m[item.domain]; !ok {
                // 注意排序函数的逻辑
                m[item.domain] = garray.NewSortedArray(100, func(v1, v2 interface{}) int {
                    item1 := v1.(*tableItem)
                    item2 := v2.(*tableItem)
                    r := 0
                    if r = strings.Compare(item1.domain, item2.domain); r == 0 {
                        if r = strings.Compare(item1.route, item2.route); r == 0 {
                            if r = strings.Compare(item1.method, item2.method); r == 0 {
                                if r = strings.Compare(item1.hook, item2.hook); r == 0 {
                                    r = item2.priority - item1.priority
                                }
                            }
                        }
                    }
                    return r
                }, false)
            }
            m[item.domain].Add(item)
        }
    }
    addr := s.config.Addr
    if s.config.HTTPSAddr != "" {
        addr += ",tls" + s.config.HTTPSAddr
    }
    for _, a := range m {
        data := make([]string, 8)
        for _, v := range a.Slice() {
            item := v.(*tableItem)
            data[0] = s.name
            data[1] = addr
            data[2] = item.domain
            data[3] = item.method
            data[4] = gconv.String(len(strings.Split(item.route, "/")) - 1 + item.priority)
            data[5] = item.route
            data[6] = item.handler
            data[7] = item.hook
            table.Append(data)
        }
    }
    table.Render()

    return buf.String()
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
                glog.Fatal(err)
            }
            // 如果所有异步的Server都已经停止，并且没有在管理操作(重启/关闭)进行中，那么主Server就可以退出了
            if serverRunning.Val() < 1 && serverProcessStatus.Val() == 0 {
                doneChan <- struct{}{}
            }
        }(v)
    }
}

// 获取当前服务器的状态
func (s *Server) Status() int {
    // 当全局运行的Web Server数量为0时表示所有Server都是停止状态
    if serverRunning.Val() == 0 {
        return SERVER_STATUS_STOPPED
    }
    // 只要有一个Server处于运行状态，那么都表示运行状态
    for _, v := range s.servers {
        if v.status == SERVER_STATUS_RUNNING {
            return SERVER_STATUS_RUNNING
        }
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
