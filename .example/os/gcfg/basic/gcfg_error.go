package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	fmt.Println(g.Config().Get("none"))
}
