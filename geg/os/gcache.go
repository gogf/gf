package main

import (
    "gitee.com/johng/gf/g/os/gcache"
    "time"
    "fmt"
    "strconv"
)

func main() {
    for i := 0; i < 10; i++ {
        gcache.Set(strconv.Itoa(i), strconv.Itoa(i), 0)
    }
    fmt.Println(gcache.Size())
    fmt.Println(gcache.Keys())
    gcache.SetCap(2)
    time.Sleep(3*time.Second)
    fmt.Println(gcache.Size())
    fmt.Println(gcache.Keys())

    return
    gcache.Set("k1", "v1", 1000)
    gcache.Set("k2", "v2", 2000)
    fmt.Println(gcache.Keys())
    fmt.Println(gcache.Values())
    fmt.Println(gcache.Size())
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