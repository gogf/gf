// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/util/gutil"
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
		if item.Handler.Type == HandlerTypeHook {
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
			// Service object.
			case HandlerTypeObject:
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
			case HandlerTypeHandler:
				m.served = true
				if m.request.IsExited() {
					break
				}
				niceCallFunc(func() {
					m.callHandlerFunc(item.Handler.Info)
				})

			// Global middleware array.
			case HandlerTypeMiddleware:
				niceCallFunc(func() {
					item.Handler.Info.Func(m.request)
				})
				// It does not continue calling next middleware after another middleware done.
				// There should be a "Next" function to be called in the middleware in order to manage the workflow.
				loop = false
			}
		}, func(exception error) {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				// It's already an error that has stack info.
				m.request.error = v
			} else {
				// Create a new error with stack info.
				// Note that there's a skip pointing the start stacktrace
				// of the real error point.
				m.request.error = gerror.WrapCodeSkip(gcode.CodeInternalError, 1, exception, "")
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
				reflect.ValueOf(m.request.Context()),
			}
			if funcInfo.Type.NumIn() == 2 {
				var (
					inputObject reflect.Value
				)
				if funcInfo.Type.In(1).Kind() == reflect.Ptr {
					inputObject = reflect.New(funcInfo.Type.In(1).Elem())
					m.request.handlerResponse.Error = m.request.Parse(inputObject.Interface())
				} else {
					inputObject = reflect.New(funcInfo.Type.In(1).Elem()).Elem()
					m.request.handlerResponse.Error = m.request.Parse(inputObject.Addr().Interface())
				}
				if m.request.handlerResponse.Error != nil {
					return
				}
				inputValues = append(inputValues, inputObject)
			}

			// Call handler with dynamic created parameter values.
			results := funcInfo.Value.Call(inputValues)
			switch len(results) {
			case 1:
				if !results[0].IsNil() {
					if err, ok := results[0].Interface().(error); ok {
						m.request.handlerResponse.Error = err
					}
				}

			case 2:
				m.request.handlerResponse.Object = results[0].Interface()
				if !results[1].IsNil() {
					if err, ok := results[1].Interface().(error); ok {
						m.request.handlerResponse.Error = err
					}
				}
			}
		}
	})
}
