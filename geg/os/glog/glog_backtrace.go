package main

import (
    "github.com/gogf/gf/g/os/glog"
    "fmt"
)

func main() {

    glog.PrintBacktrace()
    glog.New().PrintBacktrace()

    fmt.Println(glog.GetBacktrace())
    fmt.Println(glog.New().GetBacktrace())
}


