package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	fmt.Println(g.Cfg().FilePath())
	g.Cfg().Dump()
}
