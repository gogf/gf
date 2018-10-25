package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gcron"
    "gitee.com/johng/gf/g/os/glog"
    "time"
)

func main() {
    gcron.Add("0 30 * * * *", func() { glog.Println("Every hour on the half hour") })
    gcron.Add("* * * * * *",  func() { glog.Println("Every second") }, "second-cron")
    gcron.Add("@hourly",      func() { glog.Println("Every hour") })
    gcron.Add("@every 1h30m", func() { glog.Println("Every hour thirty") })
    g.Dump(gcron.Entries())

    time.Sleep(3*time.Second)

    gcron.Stop("second-cron")

    time.Sleep(3*time.Second)

    gcron.Start("second-cron")

    time.Sleep(3*time.Second)
}