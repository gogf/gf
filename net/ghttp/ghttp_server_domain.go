// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"strings"
)

// Domain is used for route register for domains.
type Domain struct {
	server  *Server             // Belonged server
	domains map[string]struct{} // Support multiple domains.
}

// Domain creates and returns a domain object for management for one or more domains.
func (s *Server) Domain(domains string) *Domain {
	d := &Domain{
		server:  s,
		domains: make(map[string]struct{}),
	}
	for _, v := range strings.Split(domains, ",") {
		d.domains[strings.TrimSpace(v)] = struct{}{}
	}
	return d
}

func (d *Domain) BindHandler(pattern string, handler interface{}) {
	for domain, _ := range d.domains {
		d.server.BindHandler(pattern+"@"+domain, handler)
	}
}

func (d *Domain) doBindHandler(ctx context.Context, in doBindHandlerInput) {
	for domain, _ := range d.domains {
		d.server.doBindHandler(ctx, doBindHandlerInput{
			Prefix:     in.Prefix,
			Pattern:    in.Pattern + "@" + domain,
			FuncInfo:   in.FuncInfo,
			Middleware: in.Middleware,
			Source:     in.Source,
		})
	}
}

func (d *Domain) BindObject(pattern string, obj interface{}, methods ...string) {
	for domain, _ := range d.domains {
		d.server.BindObject(pattern+"@"+domain, obj, methods...)
	}
}

func (d *Domain) doBindObject(ctx context.Context, in doBindObjectInput) {
	for domain, _ := range d.domains {
		d.server.doBindObject(ctx, doBindObjectInput{
			Prefix:     in.Prefix,
			Pattern:    in.Pattern + "@" + domain,
			Object:     in.Object,
			Method:     in.Method,
			Middleware: in.Middleware,
			Source:     in.Source,
		})
	}
}

func (d *Domain) BindObjectMethod(pattern string, obj interface{}, method string) {
	for domain, _ := range d.domains {
		d.server.BindObjectMethod(pattern+"@"+domain, obj, method)
	}
}

func (d *Domain) doBindObjectMethod(ctx context.Context, in doBindObjectMethodInput) {
	for domain, _ := range d.domains {
		d.server.doBindObjectMethod(ctx, doBindObjectMethodInput{
			Prefix:     in.Prefix,
			Pattern:    in.Pattern + "@" + domain,
			Object:     in.Object,
			Method:     in.Method,
			Middleware: in.Middleware,
			Source:     in.Source,
		})
	}
}

func (d *Domain) BindObjectRest(pattern string, obj interface{}) {
	for domain, _ := range d.domains {
		d.server.BindObjectRest(pattern+"@"+domain, obj)
	}
}

func (d *Domain) doBindObjectRest(ctx context.Context, in doBindObjectInput) {
	for domain, _ := range d.domains {
		d.server.doBindObjectRest(ctx, doBindObjectInput{
			Prefix:     in.Prefix,
			Pattern:    in.Pattern + "@" + domain,
			Object:     in.Object,
			Method:     in.Method,
			Middleware: in.Middleware,
			Source:     in.Source,
		})
	}
}

func (d *Domain) BindHookHandler(pattern string, hook string, handler HandlerFunc) {
	for domain, _ := range d.domains {
		d.server.BindHookHandler(pattern+"@"+domain, hook, handler)
	}
}

func (d *Domain) doBindHookHandler(ctx context.Context, in doBindHookHandlerInput) {
	for domain, _ := range d.domains {
		d.server.doBindHookHandler(ctx, doBindHookHandlerInput{
			Prefix:   in.Prefix,
			Pattern:  in.Pattern + "@" + domain,
			HookName: in.HookName,
			Handler:  in.Handler,
			Source:   in.Source,
		})
	}
}

func (d *Domain) BindHookHandlerByMap(pattern string, hookmap map[string]HandlerFunc) {
	for domain, _ := range d.domains {
		d.server.BindHookHandlerByMap(pattern+"@"+domain, hookmap)
	}
}

func (d *Domain) BindStatusHandler(status int, handler HandlerFunc) {
	for domain, _ := range d.domains {
		d.server.addStatusHandler(d.server.statusHandlerKey(status, domain), handler)
	}
}

func (d *Domain) BindStatusHandlerByMap(handlerMap map[int]HandlerFunc) {
	for k, v := range handlerMap {
		d.BindStatusHandler(k, v)
	}
}

func (d *Domain) BindMiddleware(pattern string, handlers ...HandlerFunc) {
	for domain, _ := range d.domains {
		d.server.BindMiddleware(pattern+"@"+domain, handlers...)
	}
}

func (d *Domain) BindMiddlewareDefault(handlers ...HandlerFunc) {
	for domain, _ := range d.domains {
		d.server.BindMiddleware(defaultMiddlewarePattern+"@"+domain, handlers...)
	}
}

func (d *Domain) Use(handlers ...HandlerFunc) {
	d.BindMiddlewareDefault(handlers...)
}
