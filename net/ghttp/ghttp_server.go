// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/ghttp/internal/swaggerui"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/protocol/goai"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	// Initialize the methods map.
	for _, v := range strings.Split(supportedHttpMethods, ",") {
		methodsMap[v] = struct{}{}
	}
}

// serverProcessInit initializes some process configurations, which can only be done once.
func serverProcessInit() {
	var (
		ctx = context.TODO()
	)
	if !serverProcessInitialized.Cas(false, true) {
		return
	}
	// This means it is a restart server, it should kill its parent before starting its listening,
	// to avoid duplicated port listening in two processes.
	if !genv.Get(adminActionRestartEnvKey).IsEmpty() {
		if p, err := os.FindProcess(gproc.PPid()); err == nil {
			if err = p.Kill(); err != nil {
				intlog.Error(ctx, err)
			}
			if _, err = p.Wait(); err != nil {
				intlog.Error(ctx, err)
			}
		} else {
			glog.Error(ctx, err)
		}
	}

	// Signal handler.
	go handleProcessSignal()

	// Process message handler.
	// It's enabled only graceful feature is enabled.
	if gracefulEnabled {
		intlog.Printf(ctx, "%d: graceful reload feature is enabled", gproc.Pid())
		go handleProcessMessage()
	} else {
		intlog.Printf(ctx, "%d: graceful reload feature is disabled", gproc.Pid())
	}

	// It's an ugly calling for better initializing the main package path
	// in source development environment. It is useful only be used in main goroutine.
	// It fails retrieving the main package path in asynchronous goroutines.
	gfile.MainPkgPath()
}

// GetServer creates and returns a server instance using given name and default configurations.
// Note that the parameter `name` should be unique for different servers. It returns an existing
// server instance if given `name` is already existing in the server mapping.
func GetServer(name ...interface{}) *Server {
	serverName := DefaultServerName
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
		openapi:          goai.New(),
	}
	// Initialize the server using default configurations.
	if err := s.SetConfig(NewConfig()); err != nil {
		panic(gerror.WrapCode(gcode.CodeInvalidConfiguration, err, ""))
	}
	// Record the server to internal server mapping by name.
	serverMapping.Set(serverName, s)
	// It enables OpenTelemetry for server in default.
	s.Use(internalMiddlewareServerTracing)
	return s
}

// Start starts listening on configured port.
// This function does not block the process, you can use function Wait blocking the process.
func (s *Server) Start() error {
	var (
		ctx = context.TODO()
	)

	// Swagger UI.
	if s.config.SwaggerPath != "" {
		swaggerui.Init()
		s.AddStaticPath(s.config.SwaggerPath, swaggerUIPackedPath)
		s.BindHookHandler(s.config.SwaggerPath+"/*", HookBeforeServe, s.swaggerUI)
		s.Logger().Debugf(
			ctx,
			`swagger ui is serving at address: %s%s/`,
			s.getListenAddress(),
			s.config.SwaggerPath,
		)
	}

	// OpenApi specification json producing handler.
	if s.config.OpenApiPath != "" {
		s.BindHandler(s.config.OpenApiPath, s.openapiSpec)
		s.Logger().Debugf(
			ctx,
			`openapi specification is serving at address: %s%s`,
			s.getListenAddress(),
			s.config.OpenApiPath,
		)
	} else {
		if s.config.SwaggerPath != "" {
			s.Logger().Notice(
				ctx,
				`openapi specification is disabled but swagger ui is serving, which might make no sense`,
			)
		} else {
			s.Logger().Debug(
				ctx,
				`openapi specification is disabled`,
			)
		}
	}
	// Register group routes.
	s.handlePreBindItems(ctx)

	// Server process initialization, which can only be initialized once.
	serverProcessInit()

	// Server can only be run once.
	if s.Status() == ServerStatusRunning {
		return gerror.NewCode(gcode.CodeInvalidOperation, "server is already running")
	}

	// Logging path setting check.
	if s.config.LogPath != "" && s.config.LogPath != s.config.Logger.GetPath() {
		if err := s.config.Logger.SetPath(s.config.LogPath); err != nil {
			return err
		}
	}
	// Default session storage.
	if s.config.SessionStorage == nil {
		path := ""
		if s.config.SessionPath != "" {
			path = gfile.Join(s.config.SessionPath, s.name)
			if !gfile.Exists(path) {
				if err := gfile.Mkdir(path); err != nil {
					return gerror.Wrapf(err, `mkdir failed for "%s"`, path)
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
			s.Logger().Fatalf(ctx, `%+v`, err)
		}
	}
	// Check the group routes again.
	s.handlePreBindItems(ctx)

	// If there's no route registered  and no static service enabled,
	// it then returns an error of invalid usage of server.
	if len(s.routesMap) == 0 && !s.config.FileServerEnabled {
		return gerror.NewCode(
			gcode.CodeInvalidOperation,
			`there's no route set or static feature enabled, did you forget import the router?`,
		)
	}

	// Start the HTTP server.
	reloaded := false
	fdMapStr := genv.Get(adminActionReloadEnvKey).String()
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
		gtimer.SetTimeout(ctx, time.Duration(s.config.GracefulTimeout)*time.Second, func(ctx context.Context) {
			if err := gproc.Send(gproc.PPid(), []byte("exit"), adminGProcCommGroup); err != nil {
				intlog.Error(ctx, "server error in process communication:", err)
			}
		})
	}
	s.initOpenApi()
	s.dumpRouterMap()
	return nil
}

func (s *Server) getListenAddress() string {
	var (
		array = gstr.SplitAndTrim(s.config.Address, ":")
		host  = `127.0.0.1`
		port  = 0
	)
	if len(array) > 1 {
		host = array[0]
		port = gconv.Int(array[1])
	} else {
		port = gconv.Int(array[0])
	}
	return fmt.Sprintf(`http://%s:%d`, host, port)
}

// DumpRouterMap dumps the router map to the log.
func (s *Server) dumpRouterMap() {
	var (
		ctx                          = context.TODO()
		routes                       = s.GetRoutes()
		headers                      = []string{"SERVER", "DOMAIN", "ADDRESS", "METHOD", "ROUTE", "HANDLER", "MIDDLEWARE"}
		isJustDefaultServerAndDomain = true
	)
	for _, item := range routes {
		if item.Server != DefaultServerName || item.Domain != DefaultDomainName {
			isJustDefaultServerAndDomain = false
			break
		}
	}
	if isJustDefaultServerAndDomain {
		headers = []string{"ADDRESS", "METHOD", "ROUTE", "HANDLER", "MIDDLEWARE"}
	}
	if s.config.DumpRouterMap && len(routes) > 0 {
		buffer := bytes.NewBuffer(nil)
		table := tablewriter.NewWriter(buffer)
		table.SetHeader(headers)
		table.SetRowLine(true)
		table.SetBorder(false)
		table.SetCenterSeparator("|")

		for _, item := range routes {
			var (
				data        = make([]string, 0)
				handlerName = gstr.TrimRightStr(item.Handler.Name, "-fm")
				middlewares = gstr.SplitAndTrim(item.Middleware, ",")
			)
			for k, v := range middlewares {
				middlewares[k] = gstr.TrimRightStr(v, "-fm")
			}
			item.Middleware = gstr.Join(middlewares, "\n")
			if isJustDefaultServerAndDomain {
				data = append(
					data,
					item.Address,
					item.Method,
					item.Route,
					handlerName,
					item.Middleware,
				)
			} else {
				data = append(
					data,
					item.Server,
					item.Domain,
					item.Address,
					item.Method,
					item.Route,
					handlerName,
					item.Middleware,
				)
			}
			table.Append(data)
		}
		table.Render()
		s.config.Logger.Header(false).Printf(ctx, "\n%s", buffer.String())
	}
}

// GetOpenApi returns the OpenApi specification management object of current server.
func (s *Server) GetOpenApi() *goai.OpenApiV3 {
	return s.openapi
}

// GetRoutes retrieves and returns the router array.
func (s *Server) GetRoutes() []RouterItem {
	var (
		m       = make(map[string]*garray.SortedArray)
		address = s.config.Address
	)
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
				Type:       registeredItem.Handler.Type,
				Middleware: array[1],
				Method:     array[2],
				Route:      array[3],
				Priority:   len(registeredItems) - index - 1,
				Handler:    registeredItem.Handler,
			}
			switch item.Handler.Type {
			case HandlerTypeObject, HandlerTypeHandler:
				item.IsServiceHandler = true

			case HandlerTypeMiddleware:
				item.Middleware = "GLOBAL MIDDLEWARE"
			}
			if len(item.Handler.Middleware) > 0 {
				for _, v := range item.Handler.Middleware {
					if item.Middleware != "" {
						item.Middleware += ","
					}
					item.Middleware += gdebug.FuncName(v)
				}
			}
			// If the domain does not exist in the dump map, it creates the map.
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
								if item1.Handler.Type == HandlerTypeMiddleware && item2.Handler.Type != HandlerTypeMiddleware {
									return -1
								} else if item1.Handler.Type == HandlerTypeMiddleware && item2.Handler.Type == HandlerTypeMiddleware {
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
	var (
		ctx = context.TODO()
	)
	if err := s.Start(); err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
	// Blocking using channel.
	<-s.closeChan
	// Remove plugins.
	if len(s.plugins) > 0 {
		for _, p := range s.plugins {
			intlog.Printf(ctx, `remove plugin: %s`, p.Name())
			if err := p.Remove(); err != nil {
				intlog.Errorf(ctx, "%+v", err)
			}
		}
	}
	s.Logger().Printf(ctx, "%d: all servers shutdown", gproc.Pid())
}

// Wait blocks to wait for all servers done.
// It's commonly used in multiple servers situation.
func Wait() {
	var (
		ctx = context.TODO()
	)
	<-allDoneChan
	// Remove plugins.
	serverMapping.Iterator(func(k string, v interface{}) bool {
		s := v.(*Server)
		if len(s.plugins) > 0 {
			for _, p := range s.plugins {
				intlog.Printf(ctx, `remove plugin: %s`, p.Name())
				if err := p.Remove(); err != nil {
					intlog.Error(ctx, err)
				}
			}
		}
		return true
	})
	glog.Printf(ctx, "%d: all servers shutdown", gproc.Pid())
}

// startServer starts the underlying server listening.
func (s *Server) startServer(fdMap listenerFdMap) {
	var (
		ctx          = context.TODO()
		httpsEnabled bool
	)
	// HTTPS
	if s.config.TLSConfig != nil || (s.config.HTTPSCertPath != "" && s.config.HTTPSKeyPath != "") {
		if len(s.config.HTTPSAddr) == 0 {
			if len(s.config.Address) > 0 {
				s.config.HTTPSAddr = s.config.Address
				s.config.Address = ""
			} else {
				s.config.HTTPSAddr = defaultHttpsAddr
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
				// The Windows OS does not support socket file descriptor passing
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
		s.config.Address = defaultHttpAddr
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
			// The Windows OS does not support socket file descriptor passing
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
	// Start listening asynchronously.
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
				s.Logger().Fatalf(ctx, `%+v`, err)
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
