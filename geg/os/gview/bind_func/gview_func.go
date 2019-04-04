package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gview"
)

// 用于测试的内置函数
func funcTest() string {
	return "test"
}

func main() {
	view := g.View()
	b, err := view.Parse("index.html", nil, gview.FuncMap{
		"test": funcTest,
	})
	fmt.Println(err)
	fmt.Println(string(b))
}
