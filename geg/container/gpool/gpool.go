package main

import (
    "gitee.com/johng/gf/g/container/gpool"
    "fmt"
    "time"
)

func main () {
    p := gpool.New(1000)
    fmt.Println(p.Get())
    p.Put(1)
    fmt.Println(p.Get())
    p.Put(2)
    time.Sleep(time.Second)
    fmt.Println(p.Get())
}
