package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gtime"
)


func main(){
    t1 := gfile.MTime("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/other/test.go")
    t2 := gtime.Second()
    fmt.Println(t1)
    fmt.Println(t2)
}