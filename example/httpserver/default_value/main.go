package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type GetListReq struct {
	g.Meta `path:"/" method:"get"`
	Type   string `v:"required#请选择内容模型" dc:"内容模型"`
	Page   int    `v:"min:0#分页号码错误"      dc:"分页号码" d:"1"`
	Size   int    `v:"max:50#分页数量最大50条" dc:"分页数量，最大50" d:"10"`
	Sort   int    `v:"in:0,1,2#排序类型不合法" dc:"排序类型(0:最新, 默认。1:活跃, 2:热度)"`
}
type GetListRes struct {
	Items []Item `dc:"内容列表"`
}

type Item struct {
	Id    int64  `dc:"内容ID"`
	Title string `dc:"内容标题"`
}

type Controller struct{}

func (Controller) GetList(ctx context.Context, req *GetListReq) (res *GetListRes, err error) {
	g.Log().Info(ctx, req)
	return
}

func main() {
	s := g.Server()
	s.Group("/content", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(&Controller{})
	})
	s.SetPort(8199)
	s.Run()
}
