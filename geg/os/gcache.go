package main

import (
    "gitee.com/johng/gf/g/os/gcache"
    "time"
    "fmt"
)

func main() {
    gcache.Set("k1", "v111111111111111111111111111111111111111111", 1000)
    //gcache.Set("k2", "v2", 2000)
    time.Sleep(time.Second)
    fmt.Println(gcache.Bytes())

    return
    gcache.Set("k2", "v2", 2000)
    time.Sleep(500*time.Millisecond)
    fmt.Println(gcache.Get("k1"))
    fmt.Println(gcache.Get("k2"))
    time.Sleep(400*time.Millisecond)
    fmt.Println(gcache.Get("k1"))
    fmt.Println(gcache.Get("k2"))
    time.Sleep(3000*time.Millisecond)
    fmt.Println(gcache.Get("k1"))
    fmt.Println(gcache.Get("k2"))
    time.Sleep(3000*time.Millisecond)
}