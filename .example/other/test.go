package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type T struct {
	ghttp.Plugin
}

func main() {
	s := g.Server()
	s.Plugin(new(T))
}
