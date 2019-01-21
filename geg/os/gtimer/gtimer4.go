package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtimer"
    "time"
)

func main() {
    interval := time.Second
    gtimer.AddTimes(interval, 2, func() {
        glog.Println("doing1")
    })
    gtimer.AddTimes(interval, 2, func() {
        glog.Println("doing2")
    })

    select { }
}
