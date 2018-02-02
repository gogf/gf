package main

import (
    "fmt"
    "time"
)

func main() {
    events1 := make(chan int, 100)
    events2 := make(chan int, 100)
    go func() {
        for{
            select {
            case t1 := <-events1:
                fmt.Println(t1)
            case t2 := <-events2:
                fmt.Println(t2)

            }
        }

    }()

    go func() {
        time.Sleep(2*time.Second)
        events1 <- 1
        events2 <- 2
        time.Sleep(2*time.Second)
        close(events1)
        close(events2)
    }()

    select {

    }
}