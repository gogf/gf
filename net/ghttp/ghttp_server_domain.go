// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"strings"
)

// 域名管理器对象
type Domain struct {
	s *Server         // 所属Server
	m map[string]bool // 多域名
}

// 生成一个域名对象, 参数 domains 支持给定多个域名。
func (s *Server) Domain(domains string) *Domain {
	d := &Domain{
		s: s,
		m: make(map[string]bool),
	}
	for _, v := range strings.Split(domains, ",") {
		d.m[strings.TrimSpace(v)] = true
	}
	return d
}

func (d *Domain) BindHandler(pattern string, handler HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindHandler(pattern+"@"+domain, handler)
	}
}

func (d *Domain) doBindHandler(pattern string, handler HandlerFunc, middleware []HandlerFunc) {
	for domain, _ := range d.m {
		d.s.doBindHandler(pattern+"@"+domain, handler, middleware)
	}
}

func (d *Domain) BindObject(pattern string, obj interface{}, methods ...string) {
	for domain, _ := range d.m {
		d.s.BindObject(pattern+"@"+domain, obj, methods...)
	}
}

func (d *Domain) doBindObject(pattern string, obj interface{}, methods string, middleware []HandlerFunc) {
	for domain, _ := range d.m {
		d.s.doBindObject(pattern+"@"+domain, obj, methods, middleware)
	}
}

func (d *Domain) BindObjectMethod(pattern string, obj interface{}, method string) {
	for domain, _ := range d.m {
		d.s.BindObjectMethod(pattern+"@"+domain, obj, method)
	}
}

func (d *Domain) doBindObjectMethod(pattern string, obj interface{}, method string, middleware []HandlerFunc) {
	for domain, _ := range d.m {
		d.s.doBindObjectMethod(pattern+"@"+domain, obj, method, middleware)
	}
}

func (d *Domain) BindObjectRest(pattern string, obj interface{}) {
	for domain, _ := range d.m {
		d.s.BindObjectRest(pattern+"@"+domain, obj)
	}
}

func (d *Domain) doBindObjectRest(pattern string, obj interface{}, middleware []HandlerFunc) {
	for domain, _ := range d.m {
		d.s.doBindObjectRest(pattern+"@"+domain, obj, middleware)
	}
}

func (d *Domain) BindController(pattern string, c Controller, methods ...string) {
	for domain, _ := range d.m {
		d.s.BindController(pattern+"@"+domain, c, methods...)
	}
}

func (d *Domain) doBindController(pattern string, c Controller, methods string, middleware []HandlerFunc) {
	for domain, _ := range d.m {
		d.s.doBindController(pattern+"@"+domain, c, methods, middleware)
	}
}

func (d *Domain) BindControllerMethod(pattern string, c Controller, method string) {
	for domain, _ := range d.m {
		d.s.BindControllerMethod(pattern+"@"+domain, c, method)
	}
}

func (d *Domain) doBindControllerMethod(pattern string, c Controller, method string, middleware []HandlerFunc) {
	for domain, _ := range d.m {
		d.s.doBindControllerMethod(pattern+"@"+domain, c, method, middleware)
	}
}

func (d *Domain) BindControllerRest(pattern string, c Controller) {
	for domain, _ := range d.m {
		d.s.BindControllerRest(pattern+"@"+domain, c)
	}
}

func (d *Domain) doBindControllerRest(pattern string, c Controller, middleware []HandlerFunc) {
	for domain, _ := range d.m {
		d.s.doBindControllerRest(pattern+"@"+domain, c, middleware)
	}
}

func (d *Domain) BindHookHandler(pattern string, hook string, handler HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindHookHandler(pattern+"@"+domain, hook, handler)
	}
}

func (d *Domain) BindHookHandlerByMap(pattern string, hookmap map[string]HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindHookHandlerByMap(pattern+"@"+domain, hookmap)
	}
}

func (d *Domain) BindStatusHandler(status int, handler HandlerFunc) {
	for domain, _ := range d.m {
		d.s.setStatusHandler(d.s.statusHandlerKey(status, domain), handler)
	}
}

func (d *Domain) BindStatusHandlerByMap(handlerMap map[int]HandlerFunc) {
	for k, v := range handlerMap {
		d.BindStatusHandler(k, v)
	}
}

func (d *Domain) BindMiddleware(pattern string, handlers ...HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindMiddleware(pattern+"@"+domain, handlers...)
	}
}

func (d *Domain) BindMiddlewareDefault(handlers ...HandlerFunc) {
	for domain, _ := range d.m {
		d.s.BindMiddleware(gDEFAULT_MIDDLEWARE_PATTERN+"@"+domain, handlers...)
	}
}

func (d *Domain) Use(handlers ...HandlerFunc) {
	d.BindMiddlewareDefault(handlers...)
}
