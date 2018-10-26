package main

import (
    "gitee.com/johng/gf/g/os/gcron"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    start := gtime.Second()
    gcron.Add("*/5 * * * * ?", func() {
        glog.Println(gtime.Second() - start)
    })

    select{}
}