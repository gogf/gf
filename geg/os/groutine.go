package main

import (
    "time"
    "gitee.com/johng/gf/g/os/groutine"
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func job(i int) {
    time.Sleep(2*time.Second)
    //fmt.Println("job done:", i)
}

func main() {
    for i := 0; i < 10; i++ {
        groutine.Add(func() {
            job(i)
        })
    }
    gtime.SetInterval(2*time.Second, func() bool {
        fmt.Println(groutine.Size())
        return true
    })
    time.Sleep(5000*time.Second)
}
