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
	"github.com/gogf/gf/internal/intlog"
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
	// Server wraps the http.Server and provides more feature.
	Server struct {
		name             string                           // Unique name for instance management.
		config           ServerConfig                     // Configuration.
		plugins          []Plugin                         // Plugin array.
		servers          []*gracefulServer                // Underlying http.Server array.
		serverCount      *gtype.Int                       // Underlying http.Server count.
		closeChan        chan struct{}                    // Used for underlying server closing event notification.
		serveTree        map[string]interface{}           // The route map tree.
		serveCache       *gcache.Cache                    // Server cache for internal usage.
		routesMap        map[string][]registeredRouteItem // Route map mainly for route dumps and repeated route checks.
		statusHandlerMap map[string]HandlerFunc           // Custom status handler map.
		sessionManager   *gsession.Manager                // Session manager.
	}

	// Router object.
	Router struct {
		Uri      string   // URI.
		Method   string   // HTTP method
		Domain   string   // Bound domain.
		RegRule  string   // Parsed regular expression for route matching.
		RegNames []string // Parsed router parameter names.
		Priority int      // Just for reference.
	}

	// Router item just for route dumps.
	RouterItem struct {
		Server           string       // Server name.
		Address          string       // Listening address.
		Domain           string       // Bound domain.
		Type             int          // Router type.
		Middleware       string       // Bound middleware.
		Method           string       // Handler method name.
		Route            string       // Route URI.
		Priority         int          // Just for reference.
		IsServiceHandler bool         // Is service handler.
		handler          *handlerItem // The handler.
	}

	// handlerItem is the registered handler for route handling,
	// including middleware and hook functions.
	handlerItem struct {
		itemId     int                // Unique handler item id mark.
		itemName   string             // Handler name, which is automatically retrieved from runtime stack when registered.
		itemType   int                // Handler type: object/handler/controller/middleware/hook.
		itemFunc   HandlerFunc        // Handler address.
		initFunc   HandlerFunc        // Initialization function when request enters the object(only available for object register type).
		shutFunc   HandlerFunc        // Shutdown function when request leaves out the object(only available for object register type).
		middleware []HandlerFunc      // Bound middleware array.
		ctrlInfo   *handlerController // Controller information for reflect usage.
		hookName   string             // Hook type name.
		router     *Router            // Router object.
	}

	// handlerParsedItem is the item parsed from URL.Path.
	handlerParsedItem struct {
		handler *handlerItem      // Handler information.
		values  map[string]string // Router values parsed from URL.Path.
	}

	// handlerController is the controller information used for reflect.
	handlerController struct {
		name    string       // Handler method name.
		reflect reflect.Type // Reflect type of the controller.
	}

	// registeredRouteItem stores the information of the router and is used for route map.
	registeredRouteItem struct {
		file    string       // Source file path and its line number.
		handler *handlerItem // Handler object.
	}

	// Request handler function.
	HandlerFunc = func(r *Request)

	// Listening file descriptor mapping.
	// The key is either "http" or "https" and the value is its FD.
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
	// methodsMap stores all supported HTTP method,
	// it is used for quick HTTP method searching using map.
	methodsMap = make(map[string]struct{})

	// serverMapping stores more than one server instances for current process.
	// The key is the name of the server, and the value is its instance.
	serverMapping = gmap.NewStrAnyMap(true)

	// serverRunning marks the running server count.
	// If there no successful server running or all servers shutdown, this value is 0.
	serverRunning = gtype.NewInt()

	// wsUpgrader is the default up-grader configuration for websocket.
	wsUpgrader = websocket.Upgrader{
		// It does not check the origin in default, the application can do it itself.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// allDoneChan is the event for all server have done its serving and exit.
	// It is used for process blocking purpose.
	allDoneChan = make(chan struct{}, 1000)

	// serverProcessInited is used for lazy initialization for server.
	// The process can only be initialized once.
	serverProcessInited = gtype.NewBool()

	// gracefulEnabled is used for graceful reload feature, which is false in default.
	gracefulEnabled = false
)

func init() {
	// Initialize the methods map.
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
		plugins:          make([]Plugin, 0),
		servers:          make([]*gracefulServer, 0),
		closeChan:        make(chan struct{}, 10000),
		serverCount:      gtype.NewInt(),
		statusHandlerMap: make(map[string]HandlerFunc),
		serveTree:        make(map[string]interface{}),
		serveCache:       gcache.New(),
		routesMap:        make(map[string][]registeredRouteItem),
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
	if s.config.SessionEnable {
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
	}

	// PProf feature.
	if s.config.PProfEnabled {
		s.EnablePProf(s.config.PProfPattern)
	}

	// Default HTTP handler.
	if s.config.Handler == nil {
		s.config.Handler = http.HandlerFunc(s.defaultHandler)
	}

	// Install external plugins.
	for _, p := range s.plugins {
		if err := p.Install(s); err != nil {
			s.Logger().Fatal(err)
		}
	}
	// Check the group routes again.
	s.handlePreBindItems()

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
// It's commonly used for single server situation.
func (s *Server) Run() {
	if err := s.Start(); err != nil {
		s.Logger().Fatal(err)
	}
	// Blocking using channel.
	<-s.closeChan
	// Remove plugins.
	if len(s.plugins) > 0 {
		for _, p := range s.plugins {
			intlog.Printf(`remove plugin: %s`, p.Name())
			p.Remove()
		}
	}
	s.Logger().Printf("[ghttp] %d: all servers shutdown", gproc.Pid())
}

// Wait blocks to wait for all servers done.
// It's commonly used in multiple servers situation.
func Wait() {
	<-allDoneChan
	// Remove plugins.
	serverMapping.Iterator(func(k string, v interface{}) bool {
		s := v.(*Server)
		if len(s.plugins) > 0 {
			for _, p := range s.plugins {
				intlog.Printf(`remove plugin: %s`, p.Name())
				p.Remove()
			}
		}
		return true
	})
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
