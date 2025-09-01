// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/internal/consts"
	"github.com/gogf/gf/v2/internal/reflection"
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

	// preBindItem is item for lazy registering feature of router group. preBindItem is not really registered
	// to server when route function of the group called but is lazily registered when server starts.
	preBindItem struct {
		group    *RouterGroup
		bindType string
		pattern  string
		object   any    // Can be handler, controller or object.
		params   []any  // Extra parameters for route registering depending on the type.
		source   string // Handler is a register at a certain source file path: line.
		bound    bool   // Is this item bound to server?
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
	routerGroup := &RouterGroup{
		domain: d,
		server: d.server,
		prefix: prefix,
	}
	if len(groups) > 0 {
		for _, nestedGroup := range groups {
			nestedGroup(routerGroup)
		}
	}
	return routerGroup
}

// Group creates and returns a subgroup of the current router group.
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

// Clone returns a new router group which is a clone of the current group.
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

// Bind does batch route registering feature for a router group.
func (g *RouterGroup) Bind(handlerOrObject ...any) *RouterGroup {
	var (
		ctx   = context.TODO()
		group = g.Clone()
	)
	for _, v := range handlerOrObject {
		var (
			item               = v
			originValueAndKind = reflection.OriginValueAndKind(item)
		)

		switch originValueAndKind.OriginKind {
		case reflect.Func, reflect.Struct:
			group = group.preBindToLocalArray(
				groupBindTypeHandler,
				"/",
				item,
			)

		default:
			g.server.Logger().Fatalf(
				ctx, "invalid bind parameter type: %v, should be route function or struct object",
				originValueAndKind.InputValue.Type(),
			)
		}
	}
	return group
}

// ALL register an http handler to give the route pattern and all http methods.
func (g *RouterGroup) ALL(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(
		groupBindTypeHandler,
		defaultMethod+":"+pattern,
		object,
		params...,
	)
}

// ALLMap registers http handlers for http methods using map.
func (g *RouterGroup) ALLMap(m map[string]any) {
	for pattern, object := range m {
		g.ALL(pattern, object)
	}
}

// Map registers http handlers for http methods using map.
func (g *RouterGroup) Map(m map[string]any) {
	for pattern, object := range m {
		g.preBindToLocalArray(groupBindTypeHandler, pattern, object)
	}
}

// GET registers an http handler to give the route pattern and the http method: GET.
func (g *RouterGroup) GET(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "GET:"+pattern, object, params...)
}

// PUT registers an http handler to give the route pattern and the http method: PUT.
func (g *RouterGroup) PUT(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "PUT:"+pattern, object, params...)
}

// POST registers an http handler to give the route pattern and the http method: POST.
func (g *RouterGroup) POST(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "POST:"+pattern, object, params...)
}

// DELETE registers an http handler to give the route pattern and the http method: DELETE.
func (g *RouterGroup) DELETE(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "DELETE:"+pattern, object, params...)
}

// PATCH registers an http handler to give the route pattern and the http method: PATCH.
func (g *RouterGroup) PATCH(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "PATCH:"+pattern, object, params...)
}

// HEAD registers an http handler to give the route pattern and the http method: HEAD.
func (g *RouterGroup) HEAD(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "HEAD:"+pattern, object, params...)
}

// CONNECT registers an http handler to give the route pattern and the http method: CONNECT.
func (g *RouterGroup) CONNECT(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "CONNECT:"+pattern, object, params...)
}

// OPTIONS register an http handler to give the route pattern and the http method: OPTIONS.
func (g *RouterGroup) OPTIONS(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "OPTIONS:"+pattern, object, params...)
}

// TRACE registers an http handler to give the route pattern and the http method: TRACE.
func (g *RouterGroup) TRACE(pattern string, object any, params ...any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, "TRACE:"+pattern, object, params...)
}

// REST registers an http handler to give the route pattern according to REST rule.
func (g *RouterGroup) REST(pattern string, object any) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeRest, pattern, object)
}

// Hook registers a hook to given route pattern.
func (g *RouterGroup) Hook(pattern string, hook HookName, handler HandlerFunc) *RouterGroup {
	return g.Clone().preBindToLocalArray(groupBindTypeHandler, pattern, handler, hook)
}

// Middleware binds one or more middleware to the router group.
func (g *RouterGroup) Middleware(handlers ...HandlerFunc) *RouterGroup {
	g.middleware = append(g.middleware, handlers...)
	return g
}

// preBindToLocalArray adds the route registering parameters to an internal variable array for lazily registering feature.
func (g *RouterGroup) preBindToLocalArray(bindType string, pattern string, object any, params ...any) *RouterGroup {
	_, file, line := gdebug.CallerWithFilter([]string{consts.StackFilterKeyForGoFrame})
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

// getPrefix returns the route prefix of the group, which recursively retrieves its parent's prefix.
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
			g.server.Logger().Fatalf(ctx, "invalid route pattern: %s", pattern)
		}
		// If there is already a domain, unset the domain field in the pattern.
		if g.domain != nil {
			domain = ""
		}
		if bindType == groupBindTypeRest {
			pattern = path
		} else {
			pattern = g.server.serveHandlerKey(
				method, path, domain,
			)
		}
	}
	// Filter repeated char '/'.
	pattern = gstr.Replace(pattern, "//", "/")

	// Convert params to a string array.
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
			in := doBindHandlerInput{
				Prefix:     prefix,
				Pattern:    pattern,
				FuncInfo:   funcInfo,
				Middleware: g.middleware,
				Source:     source,
			}
			if g.domain != nil {
				g.domain.doBindHandler(ctx, in)
			} else {
				g.server.doBindHandler(ctx, in)
			}
		} else {
			if len(extras) > 0 {
				if gstr.Contains(extras[0], ",") {
					in := doBindObjectInput{
						Prefix:     prefix,
						Pattern:    pattern,
						Object:     object,
						Method:     extras[0],
						Middleware: g.middleware,
						Source:     source,
					}
					if g.domain != nil {
						g.domain.doBindObject(ctx, in)
					} else {
						g.server.doBindObject(ctx, in)
					}
				} else {
					in := doBindObjectMethodInput{
						Prefix:     prefix,
						Pattern:    pattern,
						Object:     object,
						Method:     extras[0],
						Middleware: g.middleware,
						Source:     source,
					}
					if g.domain != nil {
						g.domain.doBindObjectMethod(ctx, in)
					} else {
						g.server.doBindObjectMethod(ctx, in)
					}
				}
			} else {
				in := doBindObjectInput{
					Prefix:     prefix,
					Pattern:    pattern,
					Object:     object,
					Method:     "",
					Middleware: g.middleware,
					Source:     source,
				}
				// Finally, it treats the `object` as the Object registering type.
				if g.domain != nil {
					g.domain.doBindObject(ctx, in)
				} else {
					g.server.doBindObject(ctx, in)
				}
			}
		}

	case groupBindTypeRest:
		in := doBindObjectInput{
			Prefix:     prefix,
			Pattern:    pattern,
			Object:     object,
			Method:     "",
			Middleware: g.middleware,
			Source:     source,
		}
		if g.domain != nil {
			g.domain.doBindObjectRest(ctx, in)
		} else {
			g.server.doBindObjectRest(ctx, in)
		}

	case groupBindTypeHook:
		if handler, ok := object.(HandlerFunc); ok {
			in := doBindHookHandlerInput{
				Prefix:   prefix,
				Pattern:  pattern,
				HookName: HookName(extras[0]),
				Handler:  handler,
				Source:   source,
			}
			if g.domain != nil {
				g.domain.doBindHookHandler(ctx, in)
			} else {
				g.server.doBindHookHandler(ctx, in)
			}
		} else {
			g.server.Logger().Fatalf(ctx, "invalid hook handler for pattern: %s", pattern)
		}
	}
	return g
}
