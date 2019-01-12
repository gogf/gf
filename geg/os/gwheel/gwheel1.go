package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gwheel"
    "time"
)

func main() {
    now      := time.Now()
    interval := 1400*time.Millisecond
    gwheel.Add(interval, func() {
        fmt.Println(time.Now(), time.Duration(time.Now().UnixNano() - now.UnixNano()))
        now = time.Now()
    })

    select { }
}
