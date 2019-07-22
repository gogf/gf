package main

import (
	"fmt"
	"github.com/gogf/gf/g"
)

func main() {
	fmt.Println(g.Config().Get("log-path"))
}
