package main

import (
    "gitee.com/johng/gf/g/os/gcron"
    "gitee.com/johng/gf/g/os/glog"
    "time"
)

func main() {
    gcron.AddSingleton("* * * * * *", func() {
        glog.Println("doing")
        time.Sleep(2*time.Second)
    })
    select { }
}