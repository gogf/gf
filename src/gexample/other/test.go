package main

import (
    "fmt"
    "g/os/gcache"
    "time"
)



func main() {
    gcache.Set("key", 10, 1000)
    time.Sleep(time.Second)
    fmt.Println(gcache.Get("key"))
    time.Sleep(time.Second)
    fmt.Println(gcache.Get("key"))
}