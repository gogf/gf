package main

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
)

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println(gtime.TimestampNanoStr())
	}
}
