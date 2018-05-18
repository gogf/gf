package main

import (
    "gitee.com/johng/gf/g/os/gflock"
    "fmt"
    "time"
)

func test() {
    l := gflock.New("1.lock")
    fmt.Println(l.Path())
    l.Lock()
    fmt.Println("lock 1")
    l.Lock()
    fmt.Println("lock 2")
}

func active() {
    i := 0
    for {
        time.Sleep(time.Second)
        i++
    }
}

func main() {
    go active()
    test()
}
