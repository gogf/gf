package main

import (
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
)


func main () {
    a := garray.NewSortedIntArray(0)
    a.Add(0)
    fmt.Println(a.Slice())
    fmt.Println(a.Search(0))
}
