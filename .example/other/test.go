package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	result, _ := g.View().ParseContent("姓名: ${.name}", g.Map{
		"name": "<script>alert('john');</script>",
	})
	fmt.Println(result)
}
