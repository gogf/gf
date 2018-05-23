package main

import (
    "gitee.com/johng/gf/g"
    "fmt"
)

func main() {

    c := g.Config()
    c.SetPath("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/os/gcfg")
    c.SetFileName("redis.yml")
    fmt.Println(c.Get("redis-cache"))


}
