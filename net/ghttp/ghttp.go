// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package ghttp provides powerful http server and simple client implements.
package ghttp

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gsession"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"time"
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
		statusHandlerMap map[string][]HandlerFunc         // Custom status handler map.
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

	// errorStack is the interface for Stack feature.
	errorStack interface {
		Error() string
		Stack() string
	}

	// Request handler function.
	HandlerFunc = func(r *Request)

	// Listening file descriptor mapping.
	// The key is either "http" or "https" and the value is its FD.
	listenerFdMap = map[string]string
)

const (
	HOOK_BEFORE_SERVE     = "HOOK_BEFORE_SERVE"  // Deprecated, use HookBeforeServe instead.
	HOOK_AFTER_SERVE      = "HOOK_AFTER_SERVE"   // Deprecated, use HookAfterServe instead.
	HOOK_BEFORE_OUTPUT    = "HOOK_BEFORE_OUTPUT" // Deprecated, use HookBeforeOutput instead.
	HOOK_AFTER_OUTPUT     = "HOOK_AFTER_OUTPUT"  // Deprecated, use HookAfterOutput instead.
	HookBeforeServe       = "HOOK_BEFORE_SERVE"
	HookAfterServe        = "HOOK_AFTER_SERVE"
	HookBeforeOutput      = "HOOK_BEFORE_OUTPUT"
	HookAfterOutput       = "HOOK_AFTER_OUTPUT"
	ServerStatusStopped   = 0
	ServerStatusRunning   = 1
	SupportedHttpMethods  = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
	defaultServerName     = "default"
	defaultDomainName     = "default"
	defaultMethod         = "ALL"
	handlerTypeHandler    = 1
	handlerTypeObject     = 2
	handlerTypeController = 3
	handlerTypeMiddleware = 4
	handlerTypeHook       = 5
	exceptionExit         = "exit"
	exceptionExitAll      = "exit_all"
	exceptionExitHook     = "exit_hook"
	routeCacheDuration    = time.Hour
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

	// wsUpGrader is the default up-grader configuration for websocket.
	wsUpGrader = websocket.Upgrader{
		// It does not check the origin in default, the application can do it itself.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// allDoneChan is the event for all server have done its serving and exit.
	// It is used for process blocking purpose.
	allDoneChan = make(chan struct{}, 1000)

	// serverProcessInitialized is used for lazy initialization for server.
	// The process can only be initialized once.
	serverProcessInitialized = gtype.NewBool()

	// gracefulEnabled is used for graceful reload feature, which is false in default.
	gracefulEnabled = false

	// defaultValueTags is the struct tag names for default value storing.
	defaultValueTags = []string{"d", "default"}
)
