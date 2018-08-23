package main


import (
    "fmt"
    "time"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    v := g.View()
    v.SetPath(`D:\Workspace\Go\GOPATH\src\gitee.com\johng\gf\geg\os\gview`)
    gtime.SetInterval(time.Second, func() bool {
        b, _ := v.Parse("gview.tpl", nil)
        fmt.Println(string(b))
        return true
    })
    select{}
}