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
	s.SetPort(8199)
	s.Run()
}
