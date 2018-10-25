package main

import (
    "fmt"
    "gitee.com/johng/gf/g/container/gtype"
)

func test() {
    defer fmt.Println(1)
    fmt.Println(2)
}

func main() {
    v := gtype.NewInt(1)
    fmt.Println(v.Set(2))
    fmt.Println(v.Set(2))
}
