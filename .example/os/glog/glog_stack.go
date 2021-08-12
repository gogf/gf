package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Log().PrintStack()

	fmt.Println(g.Log().GetStack())
}
