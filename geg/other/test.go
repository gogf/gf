package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "time"
    "reflect"
)

func main() {

    fmt.Println(reflect.TypeOf(gconv.Time(gtime.Second())))
    fmt.Println(time.Unix(gtime.Second(), 0).String())

}