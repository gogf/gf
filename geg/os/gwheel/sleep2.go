package main

import (
    "fmt"
    "time"
)

func main() {
    i := 0
    for {
        time.Sleep(10*time.Millisecond)
        fmt.Println(time.Now())
        i++
        if i == 100 {
            break
        }
    }
}
