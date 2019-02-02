package main

import (
    "fmt"
    "sync"
    "github.com/gogf/gf/g/os/grpool"
)

func main() {
    wg := sync.WaitGroup{}
    for i := 0; i < 10; i++ {
        wg.Add(1)
        grpool.Add(func() {
            fmt.Println(i)
            wg.Done()
        })
    }
    wg.Wait()
}
