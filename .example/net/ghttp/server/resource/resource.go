package main

import (
	_ "github.com/gogf/gf/.example/net/ghttp/server/resource/testdata"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gres"
)

func main() {
	gres.Dump()

	s := g.Server()
	s.SetIndexFolder(true)
	s.SetResource(gres.Default())
	s.SetServerRoot("/root")
	s.SetPort(8199)
	s.Run()
}
