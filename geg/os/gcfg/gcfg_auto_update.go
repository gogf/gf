package main

import (
    "fmt"
    "time"
    "gitee.com/johng/gf/g/os/gcfg"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    c := gcfg.New("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/os/gcfg")
    // 每隔1秒打印当前配置项值，用户可手动在外部修改文件内容，gcfg读取到的配置项值会即时得到更新
    gtime.SetInterval(time.Second, func() bool {
        fmt.Println(c.Get("viewpath"))
        return true
    })

    select {}
}

