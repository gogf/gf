package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gwheel"
    "time"
)

func main() {
    gwheel.Add(time.Second, func() {
        fmt.Println(time.Now().String())
    })
    select { }
}
