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
		source     string             // Source file path:line when registering.
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
		source  string       // Source file path and its line number.
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

// SetGraceful enables/disables the graceful reload feature for server,
// which is false in default.
//
// Note that this feature switch is not for single server instance but for whole process.
func SetGraceful(enabled bool) {
	gracefulEnabled = enabled
}

// serverProcessInit initializes some process configurations, which can only be done once.
func serverProcessInit() {
	if !serverProcessInited.Cas(false, true) {
		return
	}
	// This means it is a restart server, it should kill its parent before starting its listening,
	// to avoid duplicated port listening in two processes.
	if genv.Get(gADMIN_ACTION_RESTART_ENVKEY) != "" {
		if p, e := os.FindProcess(gproc.PPid()); e == nil {
			p.Kill()
			p.Wait()
		} else {
			glog.Error(e)
		}
	}

	// Signal handler.
	go handleProcessSignal()

	// Process message handler.
	// It's enabled only graceful feature is enabled.
	if gracefulEnabled {
		go handleProcessMessage()
	}

	// It's an ugly calling for better initializing the main package path
	// in source development environment. It is useful only be used in main goroutine.
	// It fails retrieving the main package path in asynchronized goroutines.
	gfile.MainPkgPath()
}

// GetServer creates and returns a server instance using given name and default configurations.
// Note that the parameter <name> should be unique for different servers. It returns an existing
// server instance if given <name> is already existing in the server mapping.
func GetServer(name ...interface{}) *Server {
	serverName := gDEFAULT_SERVER
	if len(name) > 0 && name[0] != "" {
		serverName = gconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
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
	// Initialize the server using default configurations.
	if err := s.SetConfig(Config()); err != nil {
		panic(err)
	}
	// Record the server to internal server mapping by name.
	serverMapping.Set(serverName, s)
	return s
}

// Start starts listening on configured port.
// This function does not block the process, you can use function Wait blocking the process.
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
		s.config.Handler = s
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

// startServer starts the underlying server listening.
func (s *Server) startServer(fdMap listenerFdMap) {
	var httpsEnabled bool
	// HTTPS
	if s.config.TLSConfig != nil || (s.config.HTTPSCertPath != "" && s.config.HTTPSKeyPath != "") {
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
				// The windows OS does not support socket file descriptor passing
				// from parent process.
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
	// HTTP
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
			// The windows OS does not support socket file descriptor passing
			// from parent process.
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
	// Start listening asynchronizedly.
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
			// The process exits if the server is closed with none closing error.
			if err != nil && !strings.EqualFold(http.ErrServerClosed.Error(), err.Error()) {
				s.Logger().Fatal(err)
			}
			// If all the underlying servers shutdown, the process exits.
			if s.serverCount.Add(-1) < 1 {
				s.closeChan <- struct{}{}
				if serverRunning.Add(-1) < 1 {
					serverMapping.Remove(s.name)
					allDoneChan <- struct{}{}
				}
			}
		}(v)
	}
}

// Status retrieves and returns the server status.
func (s *Server) Status() int {
	if serverRunning.Val() == 0 {
		return SERVER_STATUS_STOPPED
	}
	// If any underlying server is running, the server status is running.
	for _, v := range s.servers {
		if v.status == SERVER_STATUS_RUNNING {
			return SERVER_STATUS_RUNNING
		}
	}
	return SERVER_STATUS_STOPPED
}

// getListenerFdMap retrieves and returns the socket file descriptors.
// The key of the returned map is "http" and "https".
func (s *Server) getListenerFdMap() map[string]string {
	m := map[string]string{
		"https": "",
		"http":  "",
	}
	for _, v := range s.servers {
		str := v.itemFunc + "#" + gconv.String(v.Fd()) + ","
		if v.isHttps {
			if len(m["https"]) > 0 {
				m["https"] += ","
			}
			m["https"] += str
		} else {
			if len(m["http"]) > 0 {
				m["http"] += ","
			}
			m["http"] += str
		}
	}
	return m
}

// IsExitError checks if given error is an exit error of server.
// This is used in old version of server for custom error handler.
// Deprecated.
func IsExitError(err interface{}) bool {
	errStr := gconv.String(err)
	if strings.EqualFold(errStr, gEXCEPTION_EXIT) ||
		strings.EqualFold(errStr, gEXCEPTION_EXIT_ALL) ||
		strings.EqualFold(errStr, gEXCEPTION_EXIT_HOOK) {
		return true
	}
	return false
}
