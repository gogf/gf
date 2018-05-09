package main

import (
    "github.com/theckman/go-flock"
    "fmt"
    "time"
)

func main() {
    fileLock := flock.NewFlock("/var/lock/go-lock.lock")
    fmt.Println(fileLock.Lock())
    fmt.Println(fileLock.Lock())
    time.Sleep(1000*time.Second)
//fmt.Println(locked)
//    fmt.Println(fileLock.Locked())
//fmt.Println(err)
//    if err != nil {
//        // handle locking error
//    }
//
//    if locked {
//        // do work
//        fileLock.Unlock()
//    }
}
