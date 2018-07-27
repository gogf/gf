package main

import (
    "fmt"
    "sync"
    "time"
    "gitee.com/johng/gf/g/os/gtime"
)


func main() {
    start := gtime.Millisecond()
    wg    := sync.WaitGroup{}
    for i := 0; i < 100000; i++ {
        wg.Add(1)
        go func() {
            time.Sleep(time.Second)
            wg.Done()
        }()
    }
    wg.Wait()
    fmt.Println("time spent:", gtime.Millisecond() - start)
}
