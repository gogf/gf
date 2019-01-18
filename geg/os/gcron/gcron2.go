package main

import (
    "gitee.com/johng/gf/g/os/gcron"
    "gitee.com/johng/gf/g/os/glog"
    "time"
)

func main() {
    cron := gcron.New()
    glog.Println("start")
    cron.DelayAddOnce(1, "* * * * * *", func() {
        glog.Println("run")
    })

    time.Sleep(10*time.Second)
}