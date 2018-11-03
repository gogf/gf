package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/third/github.com/fsnotify/fsnotify"
    "os"
)

func main() {
    if w, err := fsnotify.NewWatcher(); err != nil {
        fmt.Println(err)
    } else {
        index := 0
        if array, err := gfile.ScanDir("/home/john", "*", true); err == nil {
            for _, path := range array {
                index++
                if err := w.Add(path); err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                } else {
                    fmt.Println(index)
                }
            }
        }
    }

}