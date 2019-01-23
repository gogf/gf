package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtimer"
    "time"
)

func main() {
    interval := time.Second
    gtimer.AddSingleton(interval, func() {
        glog.Println("doing")
        time.Sleep(5*time.Second)
    })

    select { }
}
