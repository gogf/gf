package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gspath"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    sp   := gspath.New()
    path := "/Users/john/Temp"
    rp, err := sp.Add(path)
    fmt.Println(err)
    fmt.Println(rp)
    fmt.Println(gtime.FuncCost(func() {
        sp.Search("1")
    }))
    fmt.Println(sp.Search("1", "index.html"))
}
