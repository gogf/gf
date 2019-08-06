// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "reflect"

// 中间件对象
type Middleware struct {
	request *Request
}

// 执行下一个请求流程处理函数
func (m *Middleware) Next() {
	item := (*handlerParsedItem)(nil)
	for {
		// 是否停止请求执行
		if m.request.IsExited() || m.request.handlerIndex >= len(m.request.handlers) {
			return
		}
		item = m.request.handlers[m.request.handlerIndex]
		m.request.handlerIndex++
		// 通过中间件模式不执行钩子函数
		if item.handler.itemType == gHANDLER_TYPE_HOOK {
			continue
		}
		// 路由参数赋值
		for k, v := range item.values {
			m.request.routerVars[k] = v
		}
		m.request.Router = item.handler.router
		// 执行函数处理
		switch item.handler.itemType {
		case gHANDLER_TYPE_CONTROLLER:
			c := reflect.New(item.handler.ctrlInfo.reflect)
			niceCallFunc(func() {
				c.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(m.request)})
			})
			if !m.request.IsExited() {
				niceCallFunc(func() {
					c.MethodByName(item.handler.ctrlInfo.name).Call(nil)
				})
			}
			if !m.request.IsExited() {
				niceCallFunc(func() {
					c.MethodByName("Shut").Call(nil)
				})
			}
		case gHANDLER_TYPE_OBJECT:
			if item.handler.initFunc != nil {
				niceCallFunc(func() {
					item.handler.initFunc(m.request)
				})
			}
			if !m.request.IsExited() {
				niceCallFunc(func() {
					item.handler.itemFunc(m.request)
				})
			}
			if !m.request.IsExited() && item.handler.shutFunc != nil {
				niceCallFunc(func() {
					item.handler.shutFunc(m.request)
				})
			}
		case gHANDLER_TYPE_MIDDLEWARE:
			niceCallFunc(func() {
				item.handler.itemFunc(m.request)
			})
		case gHANDLER_TYPE_HANDLER:
			niceCallFunc(func() {
				item.handler.itemFunc(m.request)
			})
		}
	}
}
