package main

import (
	"fmt"
	"github.com/gogf/gf/text/gstr"
)

func main() {
	filename := "FDJT02·WS·2013·DQ·D30·0002-11"
	fmt.Println(len(filename))
	fmt.Println(gstr.RuneLen(filename)) //值是29
	fmt.Println(gstr.SubStrRune(filename, 0, 26))
}
