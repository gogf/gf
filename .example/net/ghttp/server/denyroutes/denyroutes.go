package main

import "github.com/gogf/gf/v2/frame/g"

func main() {
	s := g.Server()
	s.SetDenyRoutes([]string{
		"/config*",
	})
	s.SetPort(8299)
	s.Run()
}
