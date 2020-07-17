package main

import (
	"github.com/jin502437344/gf/frame/g"
)

func main() {
	s := g.Server()
	s.EnableAdmin()
	s.SetPort(8199)
	s.Run()
}
