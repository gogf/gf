// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(
			User,
		)
	})
	oai := s.GetOpenApi()
	oai.Config.CommonResponse = ghttp.DefaultHandlerResponse{}
	oai.Config.CommonResponseDataField = "Data"
	s.SetOpenApiPath("/api")
	s.SetSwaggerPath("/swagger")
	// if api.json requires authentication, add openApiBasicAuth handler
	s.BindHookHandler(s.GetOpenApiPath(), ghttp.HookBeforeServe, openApiBasicAuth)
	s.SetPort(8199)
	s.Run()
}

func openApiBasicAuth(r *ghttp.Request) {
	if !r.BasicAuth("OpenApiAuthUserName", "OpenApiAuthPass", "Restricted") {
		r.ExitAll()
		return
	}
}
