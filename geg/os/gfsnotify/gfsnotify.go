package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfsnotify"
)

func main() {
    err := gfsnotify.Add("/home/john/Documents/temp", func(event *gfsnotify.Event) {
        fmt.Println(event)
    })
    fmt.Println(err)
    select {

    }
}