package main

import (
	_ "github.com/gogf/gf/.example/os/gres/testdata"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gres"
)

func main() {
	gres.Dump()
	g.Dump(gres.Scan("/root/image/", "*", true))
}
