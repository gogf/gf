package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
)

func main() {
    if w, err := fsnotify.NewWatcher(); err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(gtime.Now().String())
        w.Add("/tmp/test")
    }

}