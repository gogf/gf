package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "math"
)

func main() {
    fmt.Println(gtime.NewFromStr("2018-10-24 00:00:00").Nanosecond())
    fmt.Println(math.MaxInt64)
    fmt.Println(gtime.Second())
    fmt.Println(gtime.Nanosecond())
}