package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    g.TryCatch(func() {
        glog.Printfln("hello")
        g.Throw("exception")
        glog.Printfln("world")
    }, func(exception interface{}) {
        glog.Error(exception)
    })
}