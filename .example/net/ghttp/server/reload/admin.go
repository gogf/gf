package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{"Graceful": true})
	s.EnableAdmin()
	s.SetPort(8199)
	s.Run()
}
