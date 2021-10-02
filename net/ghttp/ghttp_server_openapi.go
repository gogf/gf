// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/protocol/goai"
	"github.com/gogf/gf/text/gstr"
)

func (s *Server) initOpenApi() {
	var (
		err    error
		method string
	)
	for _, item := range s.GetRoutes() {
		method = item.Method
		if gstr.Equal(method, defaultMethod) {
			method = "POST"
		}
		if item.Handler.Info.Func == nil {
			err = s.openapi.Add(goai.AddInput{
				Path:   item.Route,
				Method: method,
				Object: item.Handler.Info.Value.Interface(),
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) openapiSpecJson(r *Request) {
	err := r.Response.WriteJson(s.openapi)
	if err != nil {
		intlog.Error(r.Context(), err)
	}
}
