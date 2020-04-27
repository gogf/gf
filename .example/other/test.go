package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	type SaveReq1 struct {
		Id   uint
		Tags string
	}
	type SaveReq2 struct {
		Id   uint
		Tags []string
	}
	r1 := SaveReq1{
		Id:   1,
		Tags: "ac",
	}
	var r2 *SaveReq2
	err := gconv.Struct(r1, &r2)
	g.Dump(err)
	g.Dump(r2)
}
