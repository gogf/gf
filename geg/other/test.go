package main

import (
    "fmt"
)

func main() {
    c1 := make(chan int, 2)
    c2 := make(chan int, 5)
    c1 <- 1
    c1 <- 2
    c2  = c1
    c2 <- 3
    fmt.Println(<-c2)
    fmt.Println(<-c2)
    fmt.Println(<-c2)
}