package main

import (
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	g.Log().SetPrefix("[API]")
	g.Log().Print("hello world")
	g.Log().Error("error occurred")
}
