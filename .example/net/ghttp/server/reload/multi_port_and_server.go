package main

import (
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	s1 := g.Server("s1")
	s1.EnableAdmin()
	s1.SetPort(8100, 8200)
	s1.Start()

	s2 := g.Server("s2")
	s2.EnableAdmin()
	s2.SetPort(8300, 8400)
	s2.Start()

	g.Wait()
}
