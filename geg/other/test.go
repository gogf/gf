package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "strconv"
)

func main() {
    fmt.Println(strconv.FormatInt(gtime.Nanosecond(), 32))
    fmt.Println(gtime.Second())
    fmt.Println(gtime.Nanosecond())
}