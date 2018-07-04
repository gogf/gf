package main

import (
    "fmt"
    "gitee.com/johng/gf/g/container/garray"
)


func main () {
    a := garray.NewSortedArray(0, 0, func(v1, v2 interface{}) int {
        if v1.(int) < v2.(int) {
            return -1
        }
        if v1.(int) > v2.(int) {
            return 1
        }
        return 0
    })
    a.Add(10)
    a.Add(20)
    a.Add(30)
    a.Add(1)
    a.Add(1)
    a.Add(1)
    a.Add(1)
    fmt.Println(a.Slice())
}
