package main

import (
    "gitee.com/johng/gf/g/container/garray"
)

func main() {
    a := garray.NewSortedIntArray(0)
    a.Add(1)
    a.Remove(0)
}