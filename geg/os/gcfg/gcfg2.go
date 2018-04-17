package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gcfg"
    "time"
)

func main() {
    c := gcfg.New("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/os/gcfg")
    fmt.Println(c.GetArray("memcache"))
    time.Sleep(10*time.Second)
    // 给你10秒钟的时间修改配置，下一次读取会自动更新
    fmt.Println(c.GetArray("memcache"))
}

