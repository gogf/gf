package main

import (
	"fmt"
	"github.com/gogf/gf/util/grand"
)

func main() {
	s := "我爱GoFrame"
	for i := 0; i <= 10; i++ {
		fmt.Println(grand.Str(s, 10))
	}
}
