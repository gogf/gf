package main


import (
    "fmt"
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/frame/gins"
)

func main() {
    v := gins.View()
    v.SetPath("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/os/gview")
    gtime.SetInterval(time.Second, func() bool {
        b, _ := v.Parse("test.tpl", nil)
        fmt.Println(string(b))
        return true
    })
    select{}
}