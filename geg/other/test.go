package main

import (
    "fmt"
    "github.com/gogf/gf/g/os/gfile"
)

func main() {
    f, e := gfile.Open("/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/third")
    fmt.Println(e)
    fmt.Println(f)
}