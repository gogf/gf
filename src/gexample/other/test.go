package main

import (
    "fmt"
    "sync"
    "time"
)

type Test struct {
    Count int
    Num   int
}

func repeatedlyLock() {
    var mutex sync.Mutex
    fmt.Println("Lock the lock. (G0)")
    mutex.Lock()
    fmt.Println("The lock is locked. (G0)")
    for i := 1; i <= 3; i++ {
        go func(i int) {
            fmt.Printf("Lock the lock. (G%d)\n", i)
            mutex.Lock()
            fmt.Printf("The lock is locked. (G%d)\n", i)
        }(i)
    }
    time.Sleep(time.Second)
    fmt.Println("Unlock the lock. (G0)")
    mutex.Unlock()
    fmt.Println("The lock is unlocked. (G0)")
    time.Sleep(time.Second)
}
var mutex sync.Mutex
var t Test

func test() {
    //mutex.Lock()
    t.Num++
    //mutex.Unlock()
}

func main() {

    mutex.Lock()
    t.Count++
    test()
    fmt.Println(t)
    mutex.Unlock()



}