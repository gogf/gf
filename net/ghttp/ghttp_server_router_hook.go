// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/debug/gdebug"
	"net/http"
)

// 绑定指定的hook回调函数, pattern参数同BindHandler，支持命名路由；hook参数的值由ghttp server设定，参数不区分大小写
func (s *Server) BindHookHandler(pattern string, hook string, handler HandlerFunc) {
	s.doBindHookHandler(pattern, hook, handler, "")
}

func (s *Server) doBindHookHandler(pattern string, hook string, handler HandlerFunc, source string) {
	s.setHandler(pattern, &handlerItem{
		itemType: gHANDLER_TYPE_HOOK,
		itemName: gdebug.FuncPath(handler),
		itemFunc: handler,
		hookName: hook,
		source:   source,
	})
}

// 通过map批量绑定回调函数
func (s *Server) BindHookHandlerByMap(pattern string, hookMap map[string]HandlerFunc) {
	for k, v := range hookMap {
		s.BindHookHandler(pattern, k, v)
	}
}

// 事件回调处理，内部使用了缓存处理.
// 并按照指定hook回调函数的优先级及注册顺序进行调用
func (s *Server) callHookHandler(hook string, r *Request) {
	hookItems := r.getHookHandlers(hook)
	if len(hookItems) > 0 {
		// 备份原有的router变量
		oldRouterMap := r.routerMap
		for _, item := range hookItems {
			r.routerMap = item.values
			// 不使用hook的router对象，保留路由注册服务的router对象，不能覆盖
			// r.Router = item.handler.router
			if err := s.niceCallHookHandler(item.handler.itemFunc, r); err != nil {
				switch err {
				case gEXCEPTION_EXIT:
					break
				case gEXCEPTION_EXIT_ALL:
					fallthrough
				case gEXCEPTION_EXIT_HOOK:
					return
				default:
					r.Response.WriteStatus(http.StatusInternalServerError, err)
					panic(err)
				}
			}
		}
		// 恢复原有的router变量
		r.routerMap = oldRouterMap
	}
}

// 友好地调用方法
func (s *Server) niceCallHookHandler(f HandlerFunc, r *Request) (err interface{}) {
	defer func() {
		err = recover()
	}()
	f(r)
	return
}
