package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	type T struct {
		UpdateTime gtime.Time
	}
	t := new(T)
	gconv.Struct(g.Map{
		"UpdateTime": gtime.Now(),
	}, t)
	fmt.Println(t.UpdateTime)
}
