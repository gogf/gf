// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/debug/gdebug"
)

// BindHookHandler registers handler for specified hook.
func (s *Server) BindHookHandler(pattern string, hook string, handler HandlerFunc) {
	s.doBindHookHandler(context.TODO(), doBindHookHandlerInput{
		Prefix:   "",
		Pattern:  pattern,
		HookName: hook,
		Handler:  handler,
		Source:   "",
	})
}

// doBindHookHandlerInput is the input for BindHookHandler.
type doBindHookHandlerInput struct {
	Prefix   string
	Pattern  string
	HookName string
	Handler  HandlerFunc
	Source   string
}

// doBindHookHandler is the internal handler for BindHookHandler.
func (s *Server) doBindHookHandler(ctx context.Context, in doBindHookHandlerInput) {
	s.setHandler(
		ctx,
		setHandlerInput{
			Prefix:  in.Prefix,
			Pattern: in.Pattern,
			HandlerItem: &handlerItem{
				Type: HandlerTypeHook,
				Name: gdebug.FuncPath(in.Handler),
				Info: handlerFuncInfo{
					Func: in.Handler,
					Type: reflect.TypeOf(in.Handler),
				},
				HookName: in.HookName,
				Source:   in.Source,
			},
		},
	)
}

// BindHookHandlerByMap registers handler for specified hook.
func (s *Server) BindHookHandlerByMap(pattern string, hookMap map[string]HandlerFunc) {
	for k, v := range hookMap {
		s.BindHookHandler(pattern, k, v)
	}
}

// callHookHandler calls the hook handler by their registered sequences.
func (s *Server) callHookHandler(hook string, r *Request) {
	if !r.hasHookHandler {
		return
	}
	hookItems := r.getHookHandlers(hook)
	if len(hookItems) > 0 {
		// Backup the old router variable map.
		oldRouterMap := r.routerMap
		for _, item := range hookItems {
			r.routerMap = item.Values
			// DO NOT USE the router of the hook handler,
			// which can overwrite the router of serving handler.
			// r.Router = item.handler.router
			if err := s.niceCallHookHandler(item.Handler.Info.Func, r); err != nil {
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
	parsedItems := make([]*handlerParsedItem, 0, 4)
	for _, v := range r.handlers {
		if v.Handler.HookName != hook {
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
