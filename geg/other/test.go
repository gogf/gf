package main

import (
    "fmt"
    "math/rand"
    "sync"
)

func main() {
    p := &sync.Pool{}
    p.Put()
    for i := 0; i < 100; i++ {
        fmt.Println(rand.Intn(200))
    }
}