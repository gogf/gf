package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/container/garray"
)

func main()  {
    array := garray.NewSortedStringArray(0, false)
    array.Add("1")
    array.Add("2")
    array.Add("3")
    array.Add("4")
    array.Add("5")
    array.Add("6")
    array.Add("7")
    array.Add("8")
    array.Add("9")
    g.Dump(array.Slice())
}