package main

import (
    "fmt"
    "sync"
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/grpool"
)

func main() {
    start := gtime.Millisecond()
    wg    := sync.WaitGroup{}
    for i := 0; i < 10000000; i++ {
        wg.Add(1)
        grpool.Add(func() {
            time.Sleep(time.Millisecond)
            wg.Done()
        })
    }
    wg.Wait()
    fmt.Println(grpool.Size())
    fmt.Println("time spent:", gtime.Millisecond() - start)
}
