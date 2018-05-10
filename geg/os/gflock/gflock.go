package main

import (
    "github.com/theckman/go-flock"
    "fmt"
)

func main() {
    fileLock := flock.NewFlock("/var/lock/go-lock.lock")

    fmt.Println(fileLock.TryLock())
    fmt.Println(fileLock.TryRLock())
    //time.Sleep(1000*time.Second)
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
