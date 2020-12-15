// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
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
	"runtime"
	"strings"
	"time"

	"github.com/gogf/gf/os/gsession"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
	"github.com/olekukonko/tablewriter"
)

func init() {
	// Initialize the methods map.
	for _, v := range strings.Split(SupportedHttpMethods, ",") {
		methodsMap[v] = struct{}{}
	}
}

// SetGraceful enables/disables the graceful reload feature for server,
// which is false in default.
//
// Note that this feature switch is not for single server instance but for whole process.
// Deprecated, use configuration of ghttp.Server for controlling this feature.
func SetGraceful(enabled bool) {
	gracefulEnabled = enabled
}

// serverProcessInit initializes some process configurations, which can only be done once.
func serverProcessInit() {
	if !serverProcessInitialized.Cas(false, true) {
		return
	}
	// This means it is a restart server, it should kill its parent before starting its listening,
	// to avoid duplicated port listening in two processes.
	if genv.Get(adminActionRestartEnvKey) != "" {
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
		intlog.Printf("%d: graceful reload feature is enabled", gproc.Pid())
		go handleProcessMessage()
	} else {
		intlog.Printf("%d: graceful reload feature is disabled", gproc.Pid())
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
	serverName := defaultServerName
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
		statusHandlerMap: make(map[string][]HandlerFunc),
		serveTree:        make(map[string]interface{}),
		serveCache:       gcache.New(),
		routesMap:        make(map[string][]registeredRouteItem),
	}
	// Initialize the server using default configurations.
	if err := s.SetConfig(NewConfig()); err != nil {
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
	if s.Status() == ServerStatusRunning {
		return errors.New("[ghttp] server is already running")
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

	// If there's no route registered  and no static service enabled,
	// it then returns an error of invalid usage of server.
	if len(s.routesMap) == 0 && !s.config.FileServerEnabled {
		return errors.New(`[ghttp] there's no route set or static feature enabled, did you forget import the router?`)
	}

	// Start the HTTP server.
	reloaded := false
	fdMapStr := genv.Get(adminActionReloadEnvKey)
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
			if err := gproc.Send(gproc.PPid(), []byte("exit"), adminGProcCommGroup); err != nil {
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
			case handlerTypeController, handlerTypeObject, handlerTypeHandler:
				item.IsServiceHandler = true
			case handlerTypeMiddleware:
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
								if item1.handler.itemType == handlerTypeMiddleware && item2.handler.itemType != handlerTypeMiddleware {
									return -1
								} else if item1.handler.itemType == handlerTypeMiddleware && item2.handler.itemType == handlerTypeMiddleware {
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
		return ServerStatusStopped
	}
	// If any underlying server is running, the server status is running.
	for _, v := range s.servers {
		if v.status == ServerStatusRunning {
			return ServerStatusRunning
		}
	}
	return ServerStatusStopped
}

// getListenerFdMap retrieves and returns the socket file descriptors.
// The key of the returned map is "http" and "https".
func (s *Server) getListenerFdMap() map[string]string {
	m := map[string]string{
		"https": "",
		"http":  "",
	}
	for _, v := range s.servers {
		str := v.address + "#" + gconv.String(v.Fd()) + ","
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
	if strings.EqualFold(errStr, exceptionExit) ||
		strings.EqualFold(errStr, exceptionExitAll) ||
		strings.EqualFold(errStr, exceptionExitHook) {
		return true
	}
	return false
}
