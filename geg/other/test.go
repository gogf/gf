package main

import (
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "fmt"
)

func main() {
    s := gtime.Second()
    t := time.Unix(s, 0)
    fmt.Println(t.Format("2006-01-02 15:04:05"))
}