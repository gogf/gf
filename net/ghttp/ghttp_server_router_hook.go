// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/debug/gdebug"
	"net/http"
)

// BindHookHandler registers handler for specified hook.
func (s *Server) BindHookHandler(pattern string, hook string, handler HandlerFunc) {
	s.doBindHookHandler(pattern, hook, handler, "")
}

func (s *Server) doBindHookHandler(pattern string, hook string, handler HandlerFunc, source string) {
	s.setHandler(pattern, &handlerItem{
		itemType: handlerTypeHook,
		itemName: gdebug.FuncPath(handler),
		itemFunc: handler,
		hookName: hook,
		source:   source,
	})
}

func (s *Server) BindHookHandlerByMap(pattern string, hookMap map[string]HandlerFunc) {
	for k, v := range hookMap {
		s.BindHookHandler(pattern, k, v)
	}
}

// callHookHandler calls the hook handler by their registered sequences.
func (s *Server) callHookHandler(hook string, r *Request) {
	hookItems := r.getHookHandlers(hook)
	if len(hookItems) > 0 {
		// Backup the old router variable map.
		oldRouterMap := r.routerMap
		for _, item := range hookItems {
			r.routerMap = item.values
			// DO NOT USE the router of the hook handler,
			// which can overwrite the router of serving handler.
			// r.Router = item.handler.router
			if err := s.niceCallHookHandler(item.handler.itemFunc, r); err != nil {
				switch err {
				case exceptionExit:
					break
				case exceptionExitAll:
					fallthrough
				case exceptionExitHook:
					return
				default:
					r.Response.WriteStatus(http.StatusInternalServerError, err)
					panic(err)
				}
			}
		}
		// Restore the old router variable map.
		r.routerMap = oldRouterMap
	}
}

// getHookHandlers retrieves and returns the hook handlers of specified hook.
func (r *Request) getHookHandlers(hook string) []*handlerParsedItem {
	if !r.hasHookHandler {
		return nil
	}
	parsedItems := make([]*handlerParsedItem, 0, 4)
	for _, v := range r.handlers {
		if v.handler.hookName != hook {
			continue
		}
		item := v
		parsedItems = append(parsedItems, item)
	}
	return parsedItems
}

// niceCallHookHandler nicely calls the hook handler function,
// which means it automatically catches and returns the possible panic error to
// avoid goroutine crash.
func (s *Server) niceCallHookHandler(f HandlerFunc, r *Request) (err interface{}) {
	defer func() {
		err = recover()
	}()
	f(r)
	return
}
