// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"reflect"
	"runtime"
)

const (
	gDEFAULT_MIDDLEWARE_PATTERN = "/*"
)

// 注册中间件，绑定到指定的路由规则上，中间件参数支持多个。
func (s *Server) BindMiddleware(pattern string, handlers ...HandlerFunc) {
	for _, handler := range handlers {
		s.setHandler(pattern, &handlerItem{
			itemType: gHANDLER_TYPE_MIDDLEWARE,
			itemName: runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			itemFunc: handler,
		})
	}
}

// 注册中间件，绑定到全局路由规则("/*")上，中间件参数支持多个。
func (s *Server) BindMiddlewareDefault(handlers ...HandlerFunc) {
	for _, handler := range handlers {
		s.setHandler(gDEFAULT_MIDDLEWARE_PATTERN, &handlerItem{
			itemType: gHANDLER_TYPE_MIDDLEWARE,
			itemName: runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			itemFunc: handler,
		})
	}
}
