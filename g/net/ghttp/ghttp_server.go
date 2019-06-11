<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package ghttp

import (
<<<<<<< HEAD
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
var readyChan     = make(chan struct{}, 100000)
// Web Server已完成服务事件通道，当有事件时表示服务完成，当前进程退出
var doneChan      = make(chan struct{}, 100000)

// Web Server进程初始化
func init() {
    go func() {
        // 等待ready消息(Run方法调用)
        <- readyChan
        // 主进程只负责创建子进程
        if !gproc.IsChild() {
            sendProcessMsg(os.Getpid(), gMSG_START, nil)
        }
        // 开启进程消息监听处理
        handleProcessMsgAndSignal()

        // 服务执行完成，需要退出
        doneChan <- struct{}{}

        if !gproc.IsChild() {
            glog.Printfln("%d: all servers shutdown", gproc.Pid())
        }
    }()
=======
    "bytes"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/container/garray"
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/os/gcache"
    "github.com/gogf/gf/g/os/genv"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gproc"
    "github.com/gogf/gf/g/os/gtimer"
    "github.com/gogf/gf/g/text/gregex"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/third/github.com/gorilla/websocket"
    "github.com/gogf/gf/third/github.com/olekukonko/tablewriter"
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
        serverCount      *gtype.Int                       // 底层http.Server数量
        closeChan        chan struct{}                    // 用以关闭事件通知的通道
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
    SERVER_STATUS_STOPPED      = 0 // Server状态：停止
    SERVER_STATUS_RUNNING      = 1 // Server状态：运行
    HOOK_BEFORE_SERVE          = "BeforeServe"
    HOOK_AFTER_SERVE           = "AfterServe"
    HOOK_BEFORE_OUTPUT         = "BeforeOutput"
    HOOK_AFTER_OUTPUT          = "AfterOutput"

    // Deprecated.
    HOOK_BEFORE_CLOSE          = "BeforeClose"
    // Deprecated.
    HOOK_AFTER_CLOSE           = "AfterClose"

    HTTP_METHODS               = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
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
    // 所有支持的HTTP Method Map(初始化时自动填充),
    // 用于快速检索需要
    methodsMap       = make(map[string]struct{})

    // WebServer表，用以存储和检索名称与Server对象之间的关联关系
    serverMapping    = gmap.NewStrAnyMap()

    // 正常运行的WebServer数量，如果没有运行、失败或者全部退出，那么该值为0
    serverRunning    = gtype.NewInt()

    // WebSocket默认配置
    wsUpgrader       = websocket.Upgrader {
        // 默认允许WebSocket请求跨域，权限控制可以由业务层自己负责，灵活度更高
        CheckOrigin: func(r *http.Request) bool {
            return true
        },
    }
    // WebServer已完成服务事件通道，当有事件时表示服务完成，当前进程退出
    allDoneChan         = make(chan struct{}, 1000)

    // 用于服务进程初始化，只能初始化一次，采用“懒初始化”(在server运行时才初始化)
    serverProcessInited = gtype.NewBool()

    // 是否开启WebServer平滑重启特性, 会开启额外的本地端口监听，用于进程管理通信(默认开启)
    gracefulEnabled     = true
)

func init() {
    for _, v := range strings.Split(HTTP_METHODS, ",") {
        methodsMap[v] = struct{}{}
    }
}

// 是否开启平滑重启特性
func SetGraceful(enabled bool) {
    gracefulEnabled = enabled
}

// Web Server进程初始化.
// 注意该方法不能放置于包初始化方法init中，不使用ghttp.Server的功能便不能初始化对应的协程goroutine逻辑.
func serverProcessInit() {
    if serverProcessInited.Val() {
        return
    }
    serverProcessInited.Set(true)
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
    if gracefulEnabled {
        go handleProcessMessage()
    }

    // 是否处于开发环境，这里调用该方法初始化main包路径值，
    // 防止异步服务goroutine获取main包路径失败，
    // 该方法只有在main协程中才会执行。
    gfile.MainPkgPath()
>>>>>>> upstream/master
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
<<<<<<< HEAD
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
=======
        closeChan        : make(chan struct{}, 100),
        serverCount      : gtype.NewInt(),
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
    // 初始化时使用默认配置
    s.SetConfig(defaultServerConfig)
    // 记录到全局ServerMap中
>>>>>>> upstream/master
    serverMapping.Set(sname, s)
    return s
}

<<<<<<< HEAD
// 作为守护协程异步执行(当同一进程中存在多个Web Server时，需要采用这种方式执行)
// 需要结合Wait方式一起使用
func (s *Server) Start() error {
    // 主进程，不执行任何业务，只负责进程管理
    if !gproc.IsChild() {
        return nil
    }

    if s.status == 1 {
        return errors.New("server is already running")
    }
=======
// 作为守护协程异步执行(当同一进程中存在多个Web Server时，需要采用这种方式执行),
// 需要结合Wait方式一起使用.
func (s *Server) Start() error {
    // 服务进程初始化，只会初始化一次
    serverProcessInit()

    // 当前Web Server状态判断
    if s.Status() == SERVER_STATUS_RUNNING {
        return errors.New("server is already running")
    }

    // 没有注册任何路由，且没有开启文件服务，那么提示错误
    if len(s.routesMap) == 0 && !s.config.FileServerEnabled {
        glog.Fatal("[ghttp] no router set or static feature enabled, did you forget import the router?")
    }

>>>>>>> upstream/master
    // 底层http server配置
    if s.config.Handler == nil {
        s.config.Handler = http.HandlerFunc(s.defaultHttpHandle)
    }
<<<<<<< HEAD
    // 开启异步关闭队列处理循环
    s.startCloseQueueLoop()
    return nil
}

=======
    // 不允许访问的路由注册(使用HOOK实现)
    // TODO 去掉HOOK的实现方式
    if s.config.DenyRoutes != nil {
        for _, v := range s.config.DenyRoutes {
            s.BindHookHandler(v, HOOK_BEFORE_SERVE, func(r *Request) {
                r.Response.WriteStatus(403)
                r.ExitAll()
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
                glog.Error("[ghttp] server error in process communication:", err)
            }
        })
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
            glog.Header(false).Println(fmt.Sprintf("\n%s", s.GetRouteMap()))
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
                m[item.domain] = garray.NewSortedArraySize(100, func(v1, v2 interface{}) int {
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
        if len(addr) > 0 {
            addr += ","
        }
        addr += "tls" + s.config.HTTPSAddr
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

>>>>>>> upstream/master
// 阻塞执行监听
func (s *Server) Run() error {
    if err := s.Start(); err != nil {
        return err
    }
<<<<<<< HEAD
    // Web Server准备就绪，待执行
    readyChan <- struct{}{}
    // 阻塞等待服务执行完成
    <- doneChan
=======
    // 阻塞等待服务执行完成
    <- s.closeChan

    glog.Printf("%d: all servers shutdown", gproc.Pid())
>>>>>>> upstream/master
    return nil
}


// 阻塞等待所有Web Server停止，常用于多Web Server场景，以及需要将Web Server异步运行的场景
// 这是一个与进程相关的方法
func Wait() {
<<<<<<< HEAD
    readyChan <- struct{}{}
    <- doneChan
=======
    // 阻塞等待服务执行完成
    <- allDoneChan

    glog.Printf("%d: all servers shutdown", gproc.Pid())
>>>>>>> upstream/master
}


// 开启底层Web Server执行
func (s *Server) startServer(fdMap listenerFdMap) {
    var httpsEnabled bool
<<<<<<< HEAD
    if len(s.config.HTTPSCertPath) > 0 && len(s.config.HTTPSKeyPath) > 0 {
=======
    // 判断是否启用HTTPS
    if len(s.config.TLSConfig.Certificates) > 0 || (len(s.config.HTTPSCertPath) > 0 && len(s.config.HTTPSKeyPath) > 0) {
>>>>>>> upstream/master
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
<<<<<<< HEAD
    for _, v := range s.servers {
        go func(server *gracefulServer) {
            var err error
            if server.isHttps {
                err = server.ListenAndServeTLS(s.config.HTTPSCertPath, s.config.HTTPSKeyPath)
=======
    serverRunning.Add(1)
    for _, v := range s.servers {
        go func(server *gracefulServer) {
            s.serverCount.Add(1)
            err := (error)(nil)
            if server.isHttps {
                err = server.ListenAndServeTLS(s.config.HTTPSCertPath, s.config.HTTPSKeyPath, &s.config.TLSConfig)
>>>>>>> upstream/master
            } else {
                err = server.ListenAndServe()
            }
            // 如果非关闭错误，那么提示报错，否则认为是正常的服务关闭操作
            if err != nil && !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
<<<<<<< HEAD
                glog.Error(err)
            }
        }(v)
    }

    s.status = 1
=======
                glog.Fatal(err)
            }
            // 如果所有异步的http.Server都已经停止，那么WebServer就可以退出了
            if s.serverCount.Add(-1) < 1 {
                s.closeChan <- struct{}{}
                // 如果所有WebServer都退出，那么退出Wait等待
                if serverRunning.Add(-1) < 1 {
                    serverMapping.Remove(s.name)
                    allDoneChan <- struct{}{}
                }
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
>>>>>>> upstream/master
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
<<<<<<< HEAD

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
=======
>>>>>>> upstream/master
