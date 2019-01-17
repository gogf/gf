package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtimer"
    "time"
)

func main() {
    interval := time.Millisecond
    gtimer.AddSingleton(interval, func() {
        glog.Println("doing")
        time.Sleep(2*time.Second)
    })

    select { }
}
