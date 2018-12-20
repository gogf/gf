package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    gdb.AddDefaultConfigNode(gdb.ConfigNode{
        Type     : "mysql",
        Linkinfo : "root:12345678@tcp(127.0.0.1:3306)/test",
    })

    if r, err := g.Database().GetOne("select now() as time"); err != nil {
        glog.Error("Mysql Init Select Now : %v", err)
    } else {
        fmt.Println(r.ToMap())
    }
}
