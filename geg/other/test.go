package main

import (
    "sync"
    "fmt"
    "strconv"
    "time"
)


type LockDemo struct {
    var1 string
    var2 string
    mu1  sync.RWMutex
    mu2  sync.RWMutex
}

func main() {
    l  := LockDemo{}
    wg := sync.WaitGroup{}
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func(i int) {
            l.mu1.Lock()
            l.mu2.Lock()
            defer l.mu2.Unlock()
            defer l.mu1.Unlock()
            l.var1 = strconv.Itoa(i)
            l.var2 = strconv.Itoa(i + 1)
            wg.Done()
        }(i)
    }
    wg.Wait()
    fmt.Println(l)
}