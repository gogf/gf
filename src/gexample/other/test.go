package main

import (
    "g/os/gcache"
    "fmt"
)



func main() {
    c1 := gcache.New()
    c2 := gcache.New()
    c1.Set("a", 1, 0)
    //c2.Import(c1.Export())
    fmt.Println(c2.Get("a"))
}