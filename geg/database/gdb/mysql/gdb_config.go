package main

import (
	"fmt"
	"github.com/gogf/gf/g"
)

func main() {
	g.Config().AddPath("/home/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/geg/frame")
	if r, err := g.DB().Table("user").Where("uid=?", 1).One(); err == nil {
		fmt.Println(r["uid"].Int())
		fmt.Println(r["name"].String())
	} else {
		fmt.Println(err)
	}
}
