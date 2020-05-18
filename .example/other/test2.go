package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.Run()
}
