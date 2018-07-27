package main

import (
    "time"
    "fmt"
    "gitee.com/johng/gf/g/os/grpool"
    "gitee.com/johng/gf/g/os/gtime"
)

func job() {
    time.Sleep(1*time.Second)
}

func main() {
    pool := grpool.New(100)
    for i := 0; i < 1000; i++ {
        pool.Add(job)
    }
    fmt.Println("worker:", pool.Size())
    fmt.Println("  jobs:", pool.Jobs())
    gtime.SetInterval(time.Second, func() bool {
       fmt.Println("worker:", pool.Size())
       fmt.Println("  jobs:", pool.Jobs())
       fmt.Println()
       return true
    })

    select {}
}
