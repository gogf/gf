package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gview"
)

func main() {
	s, err := gview.ParseContent(`{{.a}}`, g.Map{
		"a": 1,
		"b": 1,
	})
	fmt.Println(err)
	fmt.Println(s)
}
