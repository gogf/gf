// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"net/http"
	"reflect"

	"github.com/gogf/gf/util/gutil"
)

// middleware is the plugin for request workflow management.
type middleware struct {
	served         bool     // Is the request served, which is used for checking response status 404.
	request        *Request // The request object pointer.
	handlerIndex   int      // Index number for executing sequence purpose for handler items.
	handlerMDIndex int      // Index number for executing sequence purpose for bound middleware of handler item.
}

// Next calls the next workflow handler.
// It's an important function controlling the workflow of the server request execution.
func (m *middleware) Next() {
	var item *handlerParsedItem
	var loop = true
	for loop {
		// Check whether the request is exited.
		if m.request.IsExited() || m.handlerIndex >= len(m.request.handlers) {
			break
		}
		item = m.request.handlers[m.handlerIndex]
		// Filter the HOOK handlers, which are designed to be called in another standalone procedure.
		if item.Handler.Type == handlerTypeHook {
			m.handlerIndex++
			continue
		}
		// Current router switching.
		m.request.Router = item.Handler.Router

		// Router values switching.
		m.request.routerMap = item.Values

		gutil.TryCatch(func() {
			// Execute bound middleware array of the item if it's not empty.
			if m.handlerMDIndex < len(item.Handler.Middleware) {
				md := item.Handler.Middleware[m.handlerMDIndex]
				m.handlerMDIndex++
				niceCallFunc(func() {
					md(m.request)
				})
				loop = false
				return
			}
			m.handlerIndex++

			switch item.Handler.Type {
			// Service controller.
			case handlerTypeController:
				m.served = true
				if m.request.IsExited() {
					break
				}
				c := reflect.New(item.Handler.CtrlInfo.Type)
				niceCallFunc(func() {
					c.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(m.request)})
				})
				if !m.request.IsExited() {
					niceCallFunc(func() {
						c.MethodByName(item.Handler.CtrlInfo.Name).Call(nil)
					})
				}
				if !m.request.IsExited() {
					niceCallFunc(func() {
						c.MethodByName("Shut").Call(nil)
					})
				}

			// Service object.
			case handlerTypeObject:
				m.served = true
				if m.request.IsExited() {
					break
				}
				if item.Handler.InitFunc != nil {
					niceCallFunc(func() {
						item.Handler.InitFunc(m.request)
					})
				}
				if !m.request.IsExited() {
					m.callHandlerFunc(item.Handler.Info)
				}
				if !m.request.IsExited() && item.Handler.ShutFunc != nil {
					niceCallFunc(func() {
						item.Handler.ShutFunc(m.request)
					})
				}

			// Service handler.
			case handlerTypeHandler:
				m.served = true
				if m.request.IsExited() {
					break
				}
				niceCallFunc(func() {
					m.callHandlerFunc(item.Handler.Info)
				})

			// Global middleware array.
			case handlerTypeMiddleware:
				niceCallFunc(func() {
					item.Handler.Info.Func(m.request)
				})
				// It does not continue calling next middleware after another middleware done.
				// There should be a "Next" function to be called in the middleware in order to manage the workflow.
				loop = false
			}
		}, func(exception error) {
			if e, ok := exception.(errorStack); ok {
				// It's already an error that has stack info.
				m.request.error = e
			} else {
				// Create a new error with stack info.
				// Note that there's a skip pointing the start stacktrace
				// of the real error point.
				m.request.error = gerror.WrapCodeSkip(gerror.CodeInternalError, 1, exception, "")
			}
			m.request.Response.WriteStatus(http.StatusInternalServerError, exception)
			loop = false
		})
	}
	// Check the http status code after all handler and middleware done.
	if m.request.IsExited() || m.handlerIndex >= len(m.request.handlers) {
		if m.request.Response.Status == 0 {
			if m.request.Middleware.served {
				m.request.Response.WriteHeader(http.StatusOK)
			} else {
				m.request.Response.WriteHeader(http.StatusNotFound)
			}
		}
	}
}

func (m *middleware) callHandlerFunc(funcInfo handlerFuncInfo) {
	niceCallFunc(func() {
		if funcInfo.Func != nil {
			funcInfo.Func(m.request)
		} else {
			var inputValues = []reflect.Value{
				reflect.ValueOf(context.WithValue(
					m.request.Context(), ctxKeyForRequest, m.request,
				)),
			}
			if funcInfo.Type.NumIn() == 2 {
				var (
					request reflect.Value
				)
				if funcInfo.Type.In(1).Kind() == reflect.Ptr {
					request = reflect.New(funcInfo.Type.In(1).Elem())
					m.request.handlerResponse.Error = m.request.Parse(request.Interface())
				} else {
					request = reflect.New(funcInfo.Type.In(1).Elem()).Elem()
					m.request.handlerResponse.Error = m.request.Parse(request.Addr().Interface())
				}
				if m.request.handlerResponse.Error != nil {
					return
				}
				inputValues = append(inputValues, request)
			}

			// Call handler with dynamic created parameter values.
			results := funcInfo.Value.Call(inputValues)
			switch len(results) {
			case 1:
				m.request.handlerResponse.Error = results[0].Interface().(error)

			case 2:
				m.request.handlerResponse.Object = results[0].Interface()
				if !results[1].IsNil() {
					if v := results[1].Interface(); v != nil {
						m.request.handlerResponse.Error = v.(error)
					}
				}
			}
		}
	})
}
