package main

import (
    "gitee.com/johng/gf/g/container/gmap"
    "fmt"
    "strings"
)

func main() {
    m := gmap.NewIntBoolMap()
    m.Set(1, true)
    fmt.Println(m.Keys())
    m.LockFunc(func(m map[int]bool) {
        m[2] = false
    })
    fmt.Println(m.Keys())
}

