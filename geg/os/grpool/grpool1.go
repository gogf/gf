package main

import (
    "time"
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/grpool"
)

func job() {
    time.Sleep(1*time.Second)
}

func main() {
    grpool.SetSize(10)
    for i := 0; i < 1000; i++ {
        grpool.Add(job)
    }
    gtime.SetInterval(2*time.Second, func() bool {
        fmt.Println("size:", grpool.Size())
        fmt.Println("jobs:", grpool.Jobs())
        return true
    })
    select {}
}
