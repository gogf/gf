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
	v.SetPath("files/template/layout2")
	s, err := v.Parse("layout.html", g.Map{
		"mainTpl": "main/main1.html",
	})
	fmt.Println(err)
	fmt.Println(s)
}
