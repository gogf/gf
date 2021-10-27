package main

import (
	"fmt"
	"github.com/gogf/gf/v2/util/guid"
)

func main() {
	for i := 0; i < 100; i++ {
		s := guid.S()
		fmt.Println(s, len(s))
	}
}
