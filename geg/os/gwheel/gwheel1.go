package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gwheel"
    "time"
)

func main() {
    fmt.Println("START:", time.Now())
    gwheel.Add(1400*time.Millisecond, func() {
        fmt.Println(time.Now())
    })

    select { }
}
