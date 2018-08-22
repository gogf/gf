package main

import (
    "gitee.com/johng/gf/g"
    "fmt"
)


func main() {
    v := g.View()
    v.AddPath("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/other")
    b, e := v.Parse("index.html")
    fmt.Println(e)
    fmt.Println(string(b))
}