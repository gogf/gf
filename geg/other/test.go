package main

import (
    "gitee.com/johng/gf/g/os/gtime"
    "fmt"
    "gitee.com/johng/gf/g/container/glist"
    "math"
)

func main() {

    t1 := gtime.Microsecond()
    c :=  make(chan struct{}, math.MaxInt64)
    c <- struct{}{}
    fmt.Println(gtime.Microsecond() - t1)

    t2 := gtime.Microsecond()
    l := glist.NewSafeList()
    l.PushBack(func() {})
    fmt.Println(gtime.Microsecond() - t2)
}