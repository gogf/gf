package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    c := g.Config()
    fmt.Println(c.GetArray("memcache"))
}

