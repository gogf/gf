package main

import (
    "time"
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/grpool"
)

func job(i int) {
    time.Sleep(2*time.Second)
    //fmt.Println("job done:", i)
}

func main() {
    for i := 0; i < 10; i++ {
        grpool.Add(func() {
            job(i)
        })
    }
    gtime.SetInterval(2*time.Second, func() bool {
        fmt.Println(grpool.Size())
        return true
    })
    time.Sleep(5000*time.Second)
}
