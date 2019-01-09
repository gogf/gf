package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("start:", time.Now())
    index  := 0
    ticker := time.NewTicker(10*time.Millisecond)
    for {
        <- ticker.C
        index++
        fmt.Println(index)
        if index == 100 {
            break
        }
    }
    fmt.Println("  end:", time.Now())
}
