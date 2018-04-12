package main

import (
    "gitee.com/johng/gf/g/encoding/gjson"
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    j := gjson.New(nil)
    t1 := gtime.Nanosecond()
    j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
    t2 := gtime.Nanosecond()
    fmt.Println(t2 - t1)

    t3 := gtime.Nanosecond()
    j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
    t4 := gtime.Nanosecond()
    fmt.Println(t4 - t3)

    t5 := gtime.Nanosecond()
    j.Get("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a")
    t6 := gtime.Nanosecond()
    fmt.Println(t6 - t5)


    j.SetViolenceCheck(false)


    t7 := gtime.Nanosecond()
    j.Set("a.a.a.a.a.a.a.a.a.a.a.a.a.a.a.a", 1)
    t8 := gtime.Nanosecond()
    fmt.Println(t8 - t7)

    t9 := gtime.Nanosecond()
    j.Get("a.a")
    t10 := gtime.Nanosecond()
    fmt.Println(t10 - t9)
}