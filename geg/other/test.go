package main

import (
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
)

func main() {
    a := garray.NewIntArray(0)
    a.Append(1)
    a.Append(2)
    a.Append(3)
    fmt.Println(a.Slice())
    a.Insert(0, 0)
    fmt.Println(a.Slice())
}

