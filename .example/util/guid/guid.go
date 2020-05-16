package main

import (
	"fmt"
	"github.com/gogf/gf/util/guid"
)

func main() {
	for i := 0; i < 1000; i++ {
		s := guid.S()
		fmt.Println(s, len(s))
	}
}
