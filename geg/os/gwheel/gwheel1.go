package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gwheel"
    "time"
)

func main() {
    _, err := gwheel.Add(time.Second, func() {
        fmt.Println(time.Now())
    })
    fmt.Println(err)
    select { }
}
