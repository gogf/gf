package main

import (
    "github.com/gogf/gf/g/os/gcron"
    "github.com/gogf/gf/g/os/glog"
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