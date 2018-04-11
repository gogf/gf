package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregx"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    t1 := gtime.Microsecond()
    for i := 0; i < 10000; i++ {
        gregx.MatchString(`([a-zA-Z]+)\^([a-zA-Z]+):(.+)@([\w\.\-]+)`, "a^b:c@d")
    }
    t2 := gtime.Microsecond()
    fmt.Println(t2 - t1)
}