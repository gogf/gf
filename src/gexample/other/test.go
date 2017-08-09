package main

import (
    "fmt"
    "sync"
    "time"
)

type T struct {
    I int
    J string
}

var m1 sync.RWMutex
var m2 sync.RWMutex

func lockPrint1(i int) {
    m1.Lock()
    m1.Lock()
    fmt.Println(i)
    m1.Unlock()
    m1.Unlock()
}

func lockPrint2(i int) {
    m2.Lock()
    fmt.Println(i)
    m2.Unlock()
}

func main() {
    go lockPrint1(1)
    //go lockPrint1(2)
    go lockPrint2(3)
    time.Sleep(2*time.Second)
}