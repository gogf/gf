package main

import (
	"fmt"
	"github.com/gogf/gf/third/golang.org/x/text/encoding/ianaindex"
)

func main() {
	e, err := ianaindex.MIB.Encoding("GB2312")
	fmt.Println(err)
	fmt.Println(e)
}
