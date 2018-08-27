package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "fmt"
)

func main() {

    glog.PrintBacktrace()
    glog.New().PrintBacktrace()

    fmt.Println(glog.GetBacktrace())
    fmt.Println(glog.New().GetBacktrace())
}


