package main

import (
	"fmt"

	"github.com/gogf/gf/util/grand"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(grand.S(16))
	}
}
