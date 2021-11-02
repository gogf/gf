package main

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	g.Log().PrintStack()

	fmt.Println(g.Log().GetStack())
}
