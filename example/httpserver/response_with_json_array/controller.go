package main

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type Req struct {
	g.Meta `path:"/user" method:"get"`
}

type Res []Item

type Item struct {
	Id   int64
	Name string
}

var (
	User = cUser{}
)

type cUser struct{}

func (c *cUser) GetList(ctx context.Context, req *Req) (res *Res, err error) {
	res = &Res{
		{Id: 1, Name: "john"},
		{Id: 2, Name: "smith"},
		{Id: 3, Name: "alice"},
		{Id: 4, Name: "katyusha"},
	}
	return
}
