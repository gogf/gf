// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/debug/gdebug"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/v2/util/gconv"
)

type (
	// RouterGroup is a group wrapping multiple routes and middleware.
	RouterGroup struct {
		parent     *RouterGroup  // Parent group.
		server     *Server       // Server.
		domain     *Domain       // Domain.
		prefix     string        // Prefix for sub-route.
		middleware []HandlerFunc // Middleware array.
	}

	// GroupItem is item for router group.
	GroupItem = []interface{}

	// preBindItem is item for lazy registering feature of router group. preBindItem is not really registered
	// to server when route function of the group called but is lazily registered when server starts.
	preBindItem struct {
		group    *RouterGroup
		bindType string
		pattern  string
		object   interface{}   // Can be handler, controller or object.
		params   []interface{} // Extra parameters for route registering depending on the type.
		source   string        // Handler is register at certain source file path:line.
		bound    bool          // Is this item bound to server.
	}
)

const (
	groupBindTypeHandler    = "HANDLER"
	groupBindTypeRest       = "REST"
	groupBindTypeHook       = "HOOK"
	groupBindTypeMiddleware = "MIDDLEWARE"
)

var (
	preBindItems = make([]*preBindItem, 0, 64)
)

// handlePreBindItems is called when server starts, which does really route registering to the server.
func (s *Server) handlePreBindItems(ctx context.Context) {
	if len(preBindItems) == 0 {
		return
	}
	for _, item := range preBindItems {
		if item.bound {
			continue
		}
		// Handle the items of current server.
		if item.group.server != nil && item.group.server != s {
			continue
		}
		if item.group.domain != nil && item.group.domain.server != s {
			continue
		}
		item.group.doBindRoutersToServer(ctx, item)
		item.bound = true
	}
}

// Group creates and returns a RouterGroup object.
func (s *Server) Group(prefix string, groups ...func(group *RouterGroup)) *RouterGroup {
	if len(prefix) > 0 && prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix == "/" {
		prefix = ""
	}
	group := &RouterGroup{
		server: s,
		prefix: prefix,
	}
	if len(groups) > 0 {
		for _, v := range groups {
			v(group)
		}
	}
	return group
}

// Group creates and returns a RouterGroup object, which is bound to a specified domain.
func (d *Domain) Group(prefix string, groups ...func(group *RouterGroup)) *RouterGroup {
	if len(prefix) > 0 && prefix[0] != '/' {
		prefix = "/" + prefix
	}
	if prefix == "/" {
		prefix = ""
	}
	group := &RouterGroup{
		domain: d,
		prefix: prefix,
	}
	if len(groups) > 0 {
		for _, v := range groups {
			v(group)
		}
	}
	return group
}

// Group creates and returns a sub-group of current router group.
func (g *RouterGroup) Group(prefix string, groups ...func(group *RouterGroup)) *RouterGroup {
	if prefix == "/" {
		prefix = ""
	}
	group := &RouterGroup{
		parent: g,
		server: g.server,
		domain: g.domain,
		prefix: prefix,
	}
	if len(g.middleware) > 0 {
		group.middleware = make([]HandlerFunc, len(g.middleware))
		copy(group.middleware, g.middleware)
	}
	if len(groups) > 0 {
		for _, v := range groups {
			v(group)
		}
	}
	return group
}

// Clone returns a new router group which is a clone of current group.
func (g *RouterGroup) Clone() *RouterGroup {
	newGroup := &RouterGroup{
		parent:     g.parent,
		server:     g.server,
		domain:     g.domain,
		prefix:     g.prefix,
		middleware: make([]HandlerFunc, len(g.middleware)),
	}
	copy(newGroup.middleware, g.middleware)
	return newGroup
}

// Bind does batch route registering feature for router group.
func (g *RouterGroup) Bind(items []GroupItem) *RouterGroup {
	var (
		ctx   = context.TODO()
		group = g.Clone()
	)
	for _, item := range items {
		if len(item) < 3 {
			g.server.Logger().Fatalf(ctx, "invalid router item: %s", item)
		}
		bindType := gstr.ToUpper(gconv.String(item[0]))
		switch bindType {
		case groupBindTypeRest:
			group.preBindToLocalArray(groupBindTypeRest, gconv.String(item[0])+":"+gconv.String(item[1]), item[2])

		case groupBindTypeMiddleware:
			group.preBindToLocalArray(groupBindTypeMiddleware, gconv.String(item[0])+":"+gconv.String(item[1]), item[2])

		default:
			if strings.EqualFold(bindType, "ALL") {
				bindType = ""
			} else {
				bindType += ":"
			}
			if len(item) > 3 {
				group.preBindToLocalArray(groupBindTypeHandler, bindType+gconv.String(item[1]), item[2], item[3])
			} else {
				group.preBindToLocalArray(groupBindTypeHandler, bindType+gconv.String(item[1]), item[2])
			}
		}
	}
	return group
}

// ALL registers a http handler to given route pattern and all http methods.
func (g *RouterGroup) ALL(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, defaultMethod+":"+pattern, object, params...)
}

// ALLMap registers http handlers for http methods using map.
func (g *RouterGroup) ALLMap(m map[string]interface{}) {
	for pattern, object := range m {
		g.ALL(pattern, object)
	}
}

// Map registers http handlers for http methods using map.
func (g *RouterGroup) Map(m map[string]interface{}) {
	for pattern, object := range m {
		g.preBindToLocalArray(groupBindTypeHandler, pattern, object)
	}
}

// GET registers a http handler to given route pattern and http method: GET.
func (g *RouterGroup) GET(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "GET:"+pattern, object, params...)
}

// PUT registers a http handler to given route pattern and http method: PUT.
func (g *RouterGroup) PUT(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "PUT:"+pattern, object, params...)
}

// POST registers a http handler to given route pattern and http method: POST.
func (g *RouterGroup) POST(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "POST:"+pattern, object, params...)
}

// DELETE registers a http handler to given route pattern and http method: DELETE.
func (g *RouterGroup) DELETE(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "DELETE:"+pattern, object, params...)
}

// PATCH registers a http handler to given route pattern and http method: PATCH.
func (g *RouterGroup) PATCH(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "PATCH:"+pattern, object, params...)
}

// HEAD registers a http handler to given route pattern and http method: HEAD.
func (g *RouterGroup) HEAD(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "HEAD:"+pattern, object, params...)
}

// CONNECT registers a http handler to given route pattern and http method: CONNECT.
func (g *RouterGroup) CONNECT(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "CONNECT:"+pattern, object, params...)
}

// OPTIONS registers a http handler to given route pattern and http method: OPTIONS.
func (g *RouterGroup) OPTIONS(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "OPTIONS:"+pattern, object, params...)
}

// TRACE registers a http handler to given route pattern and http method: TRACE.
func (g *RouterGroup) TRACE(pattern string, object interface{}, params ...interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "TRACE:"+pattern, object, params...)
}

// REST registers a http handler to given route pattern according to REST rule.
func (g *RouterGroup) REST(pattern string, object interface{}) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeRest, pattern, object)
}

// Hook registers a hook to given route pattern.
func (g *RouterGroup) Hook(pattern string, hook string, handler HandlerFunc) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, pattern, handler, hook)
}

// Middleware binds one or more middleware to the router group.
func (g *RouterGroup) Middleware(handlers ...HandlerFunc) *RouterGroup {
	g.middleware = append(g.middleware, handlers...)
	return g
}

// preBindToLocalArray adds the route registering parameters to internal variable array for lazily registering feature.
func (g *RouterGroup) preBindToLocalArray(bindType string, pattern string, object interface{}, params ...interface{}) *RouterGroup {
	_, file, line := gdebug.CallerWithFilter(stackFilterKey)
	preBindItems = append(preBindItems, &preBindItem{
		group:    g,
		bindType: bindType,
		pattern:  pattern,
		object:   object,
		params:   params,
		source:   fmt.Sprintf(`%s:%d`, file, line),
	})
	return g
}

// getPrefix returns the route prefix of the group, which recursively retrieves its parent's prefixo.
func (g *RouterGroup) getPrefix() string {
	prefix := g.prefix
	parent := g.parent
	for parent != nil {
		prefix = parent.prefix + prefix
		parent = parent.parent
	}
	return prefix
}

// doBindRoutersToServer does really register for the group.
func (g *RouterGroup) doBindRoutersToServer(ctx context.Context, item *preBindItem) *RouterGroup {
	var (
		bindType = item.bindType
		pattern  = item.pattern
		object   = item.object
		params   = item.params
		source   = item.source
	)
	prefix := g.getPrefix()
	// Route check.
	if len(prefix) > 0 {
		domain, method, path, err := g.server.parsePattern(pattern)
		if err != nil {
			g.server.Logger().Fatalf(ctx, "invalid pattern: %s", pattern)
		}
		// If there is already a domain, unset the domain field in the pattern.
		if g.domain != nil {
			domain = ""
		}
		if bindType == groupBindTypeRest {
			pattern = prefix + "/" + strings.TrimLeft(path, "/")
		} else {
			pattern = g.server.serveHandlerKey(
				method, prefix+"/"+strings.TrimLeft(path, "/"), domain,
			)
		}
	}
	// Filter repeated char '/'.
	pattern = gstr.Replace(pattern, "//", "/")

	// Convert params to string array.
	extras := gconv.Strings(params)

	// Check whether it's a hook handler.
	if _, ok := object.(HandlerFunc); ok && len(extras) > 0 {
		bindType = groupBindTypeHook
	}
	switch bindType {
	case groupBindTypeHandler:
		if reflect.ValueOf(object).Kind() == reflect.Func {
			funcInfo, err := g.server.checkAndCreateFuncInfo(object, "", "", "")
			if err != nil {
				g.server.Logger().Fatal(ctx, err.Error())
				return g
			}
			if g.server != nil {
				g.server.doBindHandler(ctx, pattern, funcInfo, g.middleware, source)
			} else {
				g.domain.doBindHandler(ctx, pattern, funcInfo, g.middleware, source)
			}
		} else {
			if len(extras) > 0 {
				if g.server != nil {
					if gstr.Contains(extras[0], ",") {
						g.server.doBindObject(
							ctx, pattern, object, extras[0], g.middleware, source,
						)
					} else {
						g.server.doBindObjectMethod(
							ctx, pattern, object, extras[0], g.middleware, source,
						)
					}
				} else {
					if gstr.Contains(extras[0], ",") {
						g.domain.doBindObject(
							ctx, pattern, object, extras[0], g.middleware, source,
						)
					} else {
						g.domain.doBindObjectMethod(
							ctx, pattern, object, extras[0], g.middleware, source,
						)
					}
				}
			} else {
				// At last, it treats the `object` as Object registering type.
				if g.server != nil {
					g.server.doBindObject(ctx, pattern, object, "", g.middleware, source)
				} else {
					g.domain.doBindObject(ctx, pattern, object, "", g.middleware, source)
				}
			}
		}

	case groupBindTypeRest:
		if g.server != nil {
			g.server.doBindObjectRest(ctx, pattern, object, g.middleware, source)
		} else {
			g.domain.doBindObjectRest(ctx, pattern, object, g.middleware, source)
		}

	case groupBindTypeHook:
		if h, ok := object.(HandlerFunc); ok {
			if g.server != nil {
				g.server.doBindHookHandler(ctx, pattern, extras[0], h, source)
			} else {
				g.domain.doBindHookHandler(ctx, pattern, extras[0], h, source)
			}
		} else {
			g.server.Logger().Fatalf(ctx, "invalid hook handler for pattern: %s", pattern)
		}
	}
	return g
}
