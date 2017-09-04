package main

import (
    "fmt"
    "g/os/gcache"
    "time"
)



func main() {
    c := gcache.New()
    c.Set("k", "v", 3)
    for {
        fmt.Println(c.Get("k"))
        time.Sleep(time.Second)
    }
}