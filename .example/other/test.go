package main

import (
	"fmt"
	"github.com/gogf/gf/text/gstr"
)

func main() {
	a := "aaaaa_post"
	b := "aaaaa_"
	c := gstr.TrimLeftStr(a, b)
	fmt.Println(c)
}
