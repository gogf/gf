package main

import (
	"fmt"
	"github.com/gogf/gf/g/util/grand"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(grand.Rand(0, 99999))
	}
}
