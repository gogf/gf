package main

import (
    "fmt"
    "time"
)

var c = make(chan int, 10)

func Chan() chan int {
    fmt.Println("yes chan")
    return c
}

func main() {
    for {
        select {
            case <- Chan():
            default:
                time.Sleep(time.Second)
        }
    }
}