package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
)


type T struct {
    name string
}

func (t *T)Test() {
    fmt.Println(t.name)
}

func main() {
    glog.Error("test")
}