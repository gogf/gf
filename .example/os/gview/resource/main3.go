package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gres"
	_ "github.com/gogf/gf/v2/os/gres/testdata"
)

func main() {
	gres.Dump()

	v := g.View()
	s, err := v.Parse("index.html")
	fmt.Println(err)
	fmt.Println(s)
}
