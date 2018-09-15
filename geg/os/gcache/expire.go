package main

import (
    "gitee.com/johng/gf/g/os/gcache"
    "fmt"
    "time"
)

func main() {
    gcache.Set("k1", "v1", 1000)
    gcache.Set("k2", "v2", 2000)
    fmt.Println(gcache.Keys())
    fmt.Println(gcache.Values())

    time.Sleep(1*time.Second)

    fmt.Println(gcache.Keys())
    fmt.Println(gcache.Values())
}