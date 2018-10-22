package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    fmt.Println(gtime.NewFromTimeStamp(gfile.MTime("/home/john/Documents/temp")).String())
}
