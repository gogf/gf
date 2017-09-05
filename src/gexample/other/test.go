package main

import (
    "fmt"
    "g/os/gcache"
)



func main() {
    c := gcache.New()
    c.Set("k", "v", 1000)
    fmt.Println(c.Get("k"))
    //c.Clear()
    fmt.Println(c.Get("k"))
    select {

    }
}