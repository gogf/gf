package main

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type HelloReq struct {
	g.Meta  `path:"/hello" tags:"Test" method:"get" description:"Hello world handler for test"`
	Content string `json:"content" in:"query"`
}

type HelloRes struct {
	Content string `json:"content" description:"Hello response content for test"`
}

// Hello is an example handler.
func Hello(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	return &HelloRes{Content: req.Content}, nil
}

func main() {
	s := g.Server()
	s.Use(
		ghttp.MiddlewareHandlerResponse,
	)
	s.BindHandler("/hello", Hello)
	s.SetOpenApiPath("/api.json")
	s.SetSwaggerPath("/swagger")
	s.SetPort(8199)
	s.Run()
}
