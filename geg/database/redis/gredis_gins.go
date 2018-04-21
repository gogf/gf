package main

import (
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/frame/gins"
)

func main() {
    gins.Config().SetPath("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame")
    gtime.SetInterval(2*time.Second, func() bool {
        redis := gins.Redis("cache")
        redis.Do("GET", "k")
        return true
    })
    select{}
}

