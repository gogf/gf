package main

import (
    "time"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    gtime.SetTimeout(3*time.Second, func() {
        glog.SetDebug(false)
    })
    for {
        glog.Debug(gtime.Datetime())
        time.Sleep(time.Second)
    }
}


