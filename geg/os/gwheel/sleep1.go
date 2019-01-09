package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    go func() {
        for {
            time.Sleep(time.Microsecond)
            go func() {
                n := 0
                for i := 0; i < 100000000; i++ {
                    n += i
                }
            }()
        }
    }()
    i := 0
    t := time.Now()
    for {
        time.Sleep(100*time.Millisecond)
        i++
        n := time.Now()
        fmt.Println(i, runtime.NumGoroutine(), n, (n.UnixNano() - t.UnixNano())/1000000)
        t = n
        if i == 100 {
            break
        }
    }
}
