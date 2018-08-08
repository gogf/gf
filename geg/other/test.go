package main

import (
	"fmt"
	"gitee.com/johng/gf/g/os/gtime"
)

func main() {
    fmt.Println(gtime.Second())
    fmt.Println(gtime.Nanosecond())
	t := gtime.Millisecond()
	for t < 1e18 {
        t *= 10
    }
	fmt.Println(t)
	fmt.Println(int64(t/1e9))
	fmt.Println(t%1e9)

	fmt.Println(gtime.NewFromTimeStamp(t).Format("Y-m-d H:i:s.u"))
}
