package main

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	t1 := gconv.Convert(1989, "Time")
	t2 := gconv.Time("2033-01-11 04:00:00 +0800 CST")
	fmt.Println(gtime.Timestamp())
	fmt.Println(t1)
	fmt.Println(t2)
}
