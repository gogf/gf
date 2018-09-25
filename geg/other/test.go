package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "time"
    "gitee.com/johng/gf/g"
)

func main() {
    g.Config().AddPath("eeee")
    g.Config().AddPath(".")
    glog.SetPath(g.Config().GetString("logPath"))
    glog.SetPath("/tmp/test-logs")
    for {
        glog.Println("1")
        time.Sleep(time.Second)
    }
}

