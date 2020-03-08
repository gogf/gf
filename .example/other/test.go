package main

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"time"
)

func main() {
	start1 := time.Now()
	end1 := start1.AddDate(0, 0, -7)
	fmt.Println(start1, end1)

	start2 := gtime.Now()
	end2 := start2.AddDate(0, 0, -7)
	fmt.Println(start2, end2)
}
