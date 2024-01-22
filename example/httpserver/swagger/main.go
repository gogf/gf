// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// HelloReq hello request
type HelloReq struct {
	g.Meta `path:"/hello" method:"get" sort:"1"`
	Name   string `v:"required" dc:"Your name"`
}

// HelloRes hello response
type HelloRes struct {
	Reply string `dc:"Reply content"`
}

// Hello Controller
type Hello struct{}

// Say function
func (Hello) Say(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	g.Log().Debugf(ctx, `receive say: %+v`, req)
	res = &HelloRes{
		Reply: fmt.Sprintf(`Hi %s`, req.Name),
	}
	return
}

func main() {
	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(Hello),
		)
	})
	s.Run()
}
