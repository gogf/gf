// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/os/gsession"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
	"github.com/gorilla/websocket"
	"github.com/olekukonko/tablewriter"
)

type (
	// Server结构体
	Server struct {
		name             string                           // 服务名称
		config           ServerConfig                     // 配置对象
		servers          []*gracefulServer                // 底层http.Server列表
		serverCount      *gtype.Int                       // 底层http.Server数量
		closeChan        chan struct{}                    // 用以关闭事件通知的通道
		servedCount      *gtype.Int                       // 已经服务的请求数(4-8字节，不考虑溢出情况)，同时作为请求ID
		serveTree        map[string]interface{}           // 所有注册的服务回调函数(路由表，树型结构，哈希表+链表优先级匹配)
		serveCache       *gcache.Cache                    // 服务注册路由内存缓存
		routesMap        map[string][]registeredRouteItem // 已经注册的路由及对应的注册方法文件地址(用以路由重复注册判断)
		statusHandlerMap map[string]HandlerFunc           // 不同状态码下的注册处理方法(例如404状态时的处理方法)
		sessionManager   *gsession.Manager                // Session管理器
	}

	// 路由对象
	Router struct {
		Uri      string   // 注册时的pattern - uri
		Method   string   // 注册时的pattern - method
		Domain   string   // 注册时的pattern - domain
		RegRule  string   // 路由规则解析后对应的正则表达式
		RegNames []string // 路由规则解析后对应的变量名称数组
		Priority int      // 优先级，用于链表排序，值越大优先级越高
	}

	// Router item just for dumping.
	RouterItem struct {
		Server           string
		Address          string
		Domain           string
		Type             int
		Middleware       string
		Method           string
		Route            string
		Priority         int
		IsServiceHandler bool
		handler          *handlerItem
	}

	// 路由函数注册信息
	handlerItem struct {
		itemId     int                // 用于标识该注册函数的唯一性ID
		itemName   string             // 注册的函数名称信息(用于路由信息打印)
		itemType   int                // 注册函数类型(对象/函数/控制器/中间件/钩子函数)
		itemFunc   HandlerFunc        // 函数内存地址(与以上两个参数二选一)
		initFunc   HandlerFunc        // 初始化请求回调函数(对象注册方式下有效)
		shutFunc   HandlerFunc        // 完成请求回调函数(对象注册方式下有效)
		middleware []HandlerFunc      // 绑定的中间件列表
		ctrlInfo   *handlerController // 控制器服务函数反射信息
		hookName   string             // 钩子类型名称(注册函数类型为钩子函数下有效)
		router     *Router            // 注册时绑定的路由对象
	}

	// 根据特定URL.Path解析后的路由检索结果项
	handlerParsedItem struct {
		handler *handlerItem      // 路由注册项
		values  map[string]string // 特定URL.Path的Router解析参数
	}

	// 控制器服务函数反射信息
	handlerController struct {
		name    string       // 方法名称
		reflect reflect.Type // 控制器类型
	}

	// 已注册的路由项
	registeredRouteItem struct {
		file    string       // 文件路径及行数地址
		handler *handlerItem // 路由注册项
	}

	// pattern与回调函数的绑定map
	handlerMap = map[string]*handlerItem

	// HTTP注册函数
	HandlerFunc = func(r *Request)

	// 文件描述符map
	listenerFdMap = map[string]string
)

const (
	SERVER_STATUS_STOPPED    = 0
	SERVER_STATUS_RUNNING    = 1
	HOOK_BEFORE_SERVE        = "HOOK_BEFORE_SERVE"
	HOOK_AFTER_SERVE         = "HOOK_AFTER_SERVE"
	HOOK_BEFORE_OUTPUT       = "HOOK_BEFORE_OUTPUT"
	HOOK_AFTER_OUTPUT        = "HOOK_AFTER_OUTPUT"
	HTTP_METHODS             = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
	gDEFAULT_SERVER          = "default"
	gDEFAULT_DOMAIN          = "default"
	gDEFAULT_METHOD          = "ALL"
	gHANDLER_TYPE_HANDLER    = 1
	gHANDLER_TYPE_OBJECT     = 2
	gHANDLER_TYPE_CONTROLLER = 3
	gHANDLER_TYPE_MIDDLEWARE = 4
	gHANDLER_TYPE_HOOK       = 5
	gEXCEPTION_EXIT          = "exit"
	gEXCEPTION_EXIT_ALL      = "exit_all"
	gEXCEPTION_EXIT_HOOK     = "exit_hook"
	gROUTE_CACHE_DURATION    = time.Hour
)

var (
	// 所有支持的HTTP Method Map(初始化时自动填充),
	// 用于快速检索需要
	methodsMap = make(map[string]struct{})

	// WebServer表，用以存储和检索名称与Server对象之间的关联关系
	serverMapping = gmap.NewStrAnyMap(true)

	// 正常运行的WebServer数量，如果没有运行、失败或者全部退出，那么该值为0
	serverRunning = gtype.NewInt()

	// WebSocket默认配置
	wsUpgrader = websocket.Upgrader{
		// 默认允许WebSocket请求跨域，权限控制可以由业务层自己负责，灵活度更高
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// WebServer已完成服务事件通道，当有事件时表示服务完成，当前进程退出
	allDoneChan = make(chan struct{}, 1000)

	// 用于服务进程初始化，只能初始化一次，采用“懒初始化”(在server运行时才初始化)
	serverProcessInited = gtype.NewBool()

	// 是否开启WebServer平滑重启特性, 会开启额外的本地端口监听，用于进程管理通信(默认开启)
	gracefulEnabled = true
)

func init() {
	for _, v := range strings.Split(HTTP_METHODS, ",") {
		methodsMap[v] = struct{}{}
	}
}

// 主要用于开发者在HTTP处理中自定义异常捕获时，判断捕获的异常是否Server抛出的自定义退出异常
func IsExitError(err interface{}) bool {
	errStr := gconv.String(err)
	if strings.EqualFold(errStr, gEXCEPTION_EXIT) ||
		strings.EqualFold(errStr, gEXCEPTION_EXIT_ALL) ||
		strings.EqualFold(errStr, gEXCEPTION_EXIT_HOOK) {
		return true
	}
	return false
}

// 是否开启平滑重启特性
func SetGraceful(enabled bool) {
	gracefulEnabled = enabled
}

// Web Server进程初始化.
// 注意该方法不能放置于包初始化方法init中，不使用ghttp.Server的功能便不能初始化对应的协程goroutine逻辑.
func serverProcessInit() {
	if !serverProcessInited.Cas(false, true) {
		return
	}
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
}

// 获取/创建一个默认配置的HTTP Server(默认监听端口是80)
// 单例模式，请保证name的唯一性
func GetServer(name ...interface{}) *Server {
	serverName := gDEFAULT_SERVER
	if len(name) > 0 && name[0] != "" {
		serverName = gconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
	c := defaultServerConfig
	s := &Server{
		name:             serverName,
		servers:          make([]*gracefulServer, 0),
		closeChan:        make(chan struct{}, 10000),
		serverCount:      gtype.NewInt(),
		statusHandlerMap: make(map[string]HandlerFunc),
		serveTree:        make(map[string]interface{}),
		serveCache:       gcache.New(),
		routesMap:        make(map[string][]registeredRouteItem),
		servedCount:      gtype.NewInt(),
	}
	// 初始化时使用默认配置
	if err := s.SetConfig(c); err != nil {
		panic(err)
	}
	// 记录到全局ServerMap中
	serverMapping.Set(serverName, s)
	return s
}

// 作为守护协程异步执行(当同一进程中存在多个Web Server时，需要采用这种方式执行),
// 需要结合Wait方式一起使用.
func (s *Server) Start() error {
	// Register group routes.
	s.handlePreBindItems()

	// Server process initialization, which can only be initialized once.
	serverProcessInit()

	// Server can only be run once.
	if s.Status() == SERVER_STATUS_RUNNING {
		return errors.New("[ghttp] server is already running")
	}

	// If there's no route registered  and no static service enabled,
	// it then returns an error of invalid usage of server.
	if len(s.routesMap) == 0 && !s.config.FileServerEnabled {
		return errors.New(`[ghttp] there's no route set or static feature enabled, did you forget import the router?`)
	}
	// Logging path setting check.
	if s.config.LogPath != "" {
		if err := s.config.Logger.SetPath(s.config.LogPath); err != nil {
			return errors.New(fmt.Sprintf("[ghttp] set log path '%s' error: %v", s.config.LogPath, err))
		}
	}
	// Default session storage.
	if s.config.SessionStorage == nil {
		path := ""
		if s.config.SessionPath != "" {
			path = gfile.Join(s.config.SessionPath, s.name)
			if !gfile.Exists(path) {
				if err := gfile.Mkdir(path); err != nil {
					return errors.New(fmt.Sprintf("[ghttp] mkdir failed for '%s': %v", path, err))
				}
			}
		}
		s.config.SessionStorage = gsession.NewStorageFile(path)
	}
	// Initialize session manager when start running.
	s.sessionManager = gsession.New(
		s.config.SessionMaxAge,
		s.config.SessionStorage,
	)

	// PProf feature.
	if s.config.PProfEnabled {
		s.EnablePProf(s.config.PProfPattern)
	}

	// Default HTTP handler.
	if s.config.Handler == nil {
		s.config.Handler = http.HandlerFunc(s.defaultHandler)
	}

	// Start the HTTP server.
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

	// If this is a child process, it then notifies its parent exit.
	if gproc.IsChild() {
		gtimer.SetTimeout(2*time.Second, func() {
			if err := gproc.Send(gproc.PPid(), []byte("exit"), gADMIN_GPROC_COMM_GROUP); err != nil {
				//glog.Error("[ghttp] server error in process communication:", err)
			}
		})
	}

	s.dumpRouterMap()
	return nil
}

// DumpRouterMap dumps the router map to the log.
func (s *Server) dumpRouterMap() {
	if s.config.DumpRouterMap && len(s.routesMap) > 0 {
		buffer := bytes.NewBuffer(nil)
		table := tablewriter.NewWriter(buffer)
		table.SetHeader([]string{"SERVER", "DOMAIN", "ADDRESS", "METHOD", "ROUTE", "HANDLER", "MIDDLEWARE"})
		table.SetRowLine(true)
		table.SetBorder(false)
		table.SetCenterSeparator("|")

		for _, item := range s.GetRouterArray() {
			data := make([]string, 7)
			data[0] = item.Server
			data[1] = item.Domain
			data[2] = item.Address
			data[3] = item.Method
			data[4] = item.Route
			data[5] = item.handler.itemName
			data[6] = item.Middleware
			table.Append(data)
		}
		table.Render()
		s.config.Logger.Header(false).Printf("\n%s", buffer.String())
	}
}

// GetRouterArray retrieves and returns the router array.
// The key of the returned map is the domain of the server.
func (s *Server) GetRouterArray() []RouterItem {
	m := make(map[string]*garray.SortedArray)
	address := s.config.Address
	if s.config.HTTPSAddr != "" {
		if len(address) > 0 {
			address += ","
		}
		address += "tls" + s.config.HTTPSAddr
	}
	for k, registeredItems := range s.routesMap {
		array, _ := gregex.MatchString(`(.*?)%([A-Z]+):(.+)@(.+)`, k)
		for index, registeredItem := range registeredItems {
			item := RouterItem{
				Server:     s.name,
				Address:    address,
				Domain:     array[4],
				Type:       registeredItem.handler.itemType,
				Middleware: array[1],
				Method:     array[2],
				Route:      array[3],
				Priority:   len(registeredItems) - index - 1,
				handler:    registeredItem.handler,
			}
			switch item.handler.itemType {
			case gHANDLER_TYPE_CONTROLLER, gHANDLER_TYPE_OBJECT, gHANDLER_TYPE_HANDLER:
				item.IsServiceHandler = true
			case gHANDLER_TYPE_MIDDLEWARE:
				item.Middleware = "GLOBAL MIDDLEWARE"
			}
			if len(item.handler.middleware) > 0 {
				for _, v := range item.handler.middleware {
					if item.Middleware != "" {
						item.Middleware += ","
					}
					item.Middleware += gdebug.FuncName(v)
				}
			}
			// If the domain does not exist in the dump map, it create the map.
			// The value of the map is a custom sorted array.
			if _, ok := m[item.Domain]; !ok {
				// Sort in ASC order.
				m[item.Domain] = garray.NewSortedArray(func(v1, v2 interface{}) int {
					item1 := v1.(RouterItem)
					item2 := v2.(RouterItem)
					r := 0
					if r = strings.Compare(item1.Domain, item2.Domain); r == 0 {
						if r = strings.Compare(item1.Route, item2.Route); r == 0 {
							if r = strings.Compare(item1.Method, item2.Method); r == 0 {
								if item1.handler.itemType == gHANDLER_TYPE_MIDDLEWARE && item2.handler.itemType != gHANDLER_TYPE_MIDDLEWARE {
									return -1
								} else if item1.handler.itemType == gHANDLER_TYPE_MIDDLEWARE && item2.handler.itemType == gHANDLER_TYPE_MIDDLEWARE {
									return 1
								} else if r = strings.Compare(item1.Middleware, item2.Middleware); r == 0 {
									r = item2.Priority - item1.Priority
								}
							}
						}
					}
					return r
				})
			}
			m[item.Domain].Add(item)
		}
	}
	routerArray := make([]RouterItem, 0, 128)
	for _, array := range m {
		for _, v := range array.Slice() {
			routerArray = append(routerArray, v.(RouterItem))
		}
	}
	return routerArray
}

// Run starts server listening in blocking way.
func (s *Server) Run() {
	if err := s.Start(); err != nil {
		s.Logger().Fatal(err)
	}

	// Blocking using channel.
	<-s.closeChan

	s.Logger().Printf("[ghttp] %d: all servers shutdown", gproc.Pid())
}

// Wait blocks to wait for all servers done.
// It's commonly used in multiple servers situation.
func Wait() {
	<-allDoneChan

	glog.Printf("[ghttp] %d: all servers shutdown", gproc.Pid())
}

// 开启底层Web Server执行
func (s *Server) startServer(fdMap listenerFdMap) {
	var httpsEnabled bool
	// 判断是否启用HTTPS
	if s.config.TLSConfig != nil || (s.config.HTTPSCertPath != "" && s.config.HTTPSKeyPath != "") {
		// ================
		// HTTPS
		// ================
		if len(s.config.HTTPSAddr) == 0 {
			if len(s.config.Address) > 0 {
				s.config.HTTPSAddr = s.config.Address
				s.config.Address = ""
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
			fd := 0
			itemFunc := v
			array := strings.Split(v, "#")
			if len(array) > 1 {
				itemFunc = array[0]
				// windows系统不支持文件描述符传递socket通信平滑交接，因此只能完整重启
				if runtime.GOOS != "windows" {
					fd = gconv.Int(array[1])
				}
			}
			if fd > 0 {
				s.servers = append(s.servers, s.newGracefulServer(itemFunc, fd))
			} else {
				s.servers = append(s.servers, s.newGracefulServer(itemFunc))
			}
			s.servers[len(s.servers)-1].isHttps = true
		}
	}
	// ================
	// HTTP
	// ================
	// 当HTTPS服务未启用时，默认HTTP地址才会生效
	if !httpsEnabled && len(s.config.Address) == 0 {
		s.config.Address = gDEFAULT_HTTP_ADDR
	}
	var array []string
	if v, ok := fdMap["http"]; ok && len(v) > 0 {
		array = strings.Split(v, ",")
	} else {
		array = strings.Split(s.config.Address, ",")
	}
	for _, v := range array {
		if len(v) == 0 {
			continue
		}
		fd := 0
		itemFunc := v
		array := strings.Split(v, "#")
		if len(array) > 1 {
			itemFunc = array[0]
			// windows系统不支持文件描述符传递socket通信平滑交接，因此只能完整重启
			if runtime.GOOS != "windows" {
				fd = gconv.Int(array[1])
			}
		}
		if fd > 0 {
			s.servers = append(s.servers, s.newGracefulServer(itemFunc, fd))
		} else {
			s.servers = append(s.servers, s.newGracefulServer(itemFunc))
		}
	}
	// 开始执行异步监听
	serverRunning.Add(1)
	for _, v := range s.servers {
		go func(server *gracefulServer) {
			s.serverCount.Add(1)
			err := (error)(nil)
			if server.isHttps {
				err = server.ListenAndServeTLS(s.config.HTTPSCertPath, s.config.HTTPSKeyPath, s.config.TLSConfig)
			} else {
				err = server.ListenAndServe()
			}
			// 如果非关闭错误，那么提示报错，否则认为是正常的服务关闭操作
			if err != nil && !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
				s.Logger().Fatal(err)
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
}

// 获取当前监听的文件描述符信息，构造成map返回
func (s *Server) getListenerFdMap() map[string]string {
	m := map[string]string{
		"https": "",
		"http":  "",
	}
	// s.servers是从HTTPS到HTTP优先级遍历，解析的时候也应当按照这个顺序读取fd
	for _, v := range s.servers {
		str := v.itemFunc + "#" + gconv.String(v.Fd()) + ","
		if v.isHttps {
			m["https"] += str
		} else {
			m["http"] += str
		}
	}
	// 去掉末尾的","号
	if len(m["https"]) > 0 {
		m["https"] = m["https"][0 : len(m["https"])-1]
	}
	if len(m["http"]) > 0 {
		m["http"] = m["http"][0 : len(m["http"])-1]
	}

	return m
}
