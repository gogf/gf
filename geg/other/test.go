package main

import (
    "gitee.com/johng/gf/g/container/gpool"
    "fmt"
    "time"
)

func main() {
    p := gpool.New(1000)
    for i := 0 ; i < 100; i++ {
        p.Put(i)
    }
    fmt.Println(p.Size())
    time.Sleep(2*time.Second)
    fmt.Println(p.Size())
}