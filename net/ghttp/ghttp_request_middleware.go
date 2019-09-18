// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"
	"reflect"

	"github.com/gogf/gf/errors/gerror"

	"github.com/gogf/gf/util/gutil"
)

// 中间件对象
type Middleware struct {
	served  bool     // 是否带有请求服务函数，用以识别是否404
	request *Request // 请求对象
}

// 执行下一个请求流程处理函数
func (m *Middleware) Next() {
	item := (*handlerParsedItem)(nil)
	loop := true
	for loop {
		// 是否停止请求执行
		if m.request.IsExited() || m.request.handlerIndex >= len(m.request.handlers) {
			return
		}
		item = m.request.handlers[m.request.handlerIndex]
		m.request.handlerIndex++
		// 中间件执行时不执行钩子函数，由另外的逻辑进行控制
		if item.handler.itemType == gHANDLER_TYPE_HOOK {
			continue
		}
		// 路由参数赋值
		for k, v := range item.values {
			m.request.routerMap[k] = v
		}
		m.request.Router = item.handler.router
		// 执行函数处理
		gutil.TryCatch(func() {
			switch item.handler.itemType {
			case gHANDLER_TYPE_CONTROLLER:
				m.served = true
				if m.request.IsExited() {
					break
				}
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
				m.served = true
				if m.request.IsExited() {
					break
				}
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

			case gHANDLER_TYPE_HANDLER:
				m.served = true
				if m.request.IsExited() {
					break
				}
				niceCallFunc(func() {
					item.handler.itemFunc(m.request)
				})

			case gHANDLER_TYPE_MIDDLEWARE:
				niceCallFunc(func() {
					item.handler.itemFunc(m.request)
				})
				// 中间件默认不会进一步执行，
				// 需要内部调用Next方法决定是否进一步执行，以便于请求流程控制。
				loop = false
			}
		}, func(exception interface{}) {
			m.request.error = gerror.Newf("%v", exception)
			m.request.Response.WriteStatus(http.StatusInternalServerError, exception)
		})
	}
}
