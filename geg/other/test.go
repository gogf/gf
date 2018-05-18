package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "time"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    file, err := gfile.Open("/home/john/Documents/temp.txt")
    fmt.Println(err)
    gtime.SetInterval(time.Second, func() bool {
        if s, e := file.Stat(); e == nil {
            fmt.Println(s.ModTime().Unix())
            fmt.Println(gfile.MTime("/home/john/Documents/temp.txt"))
        }
        return true
    })

    time.Sleep(time.Hour)

}
