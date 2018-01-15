package main

import (
    "time"
    "gitee.com/johng/gf/g/os/groutine"
    "fmt"
)

func job() {
    time.Sleep(3*time.Second)
    fmt.Println("job done")
}

func main() {
    p := groutine.New()
    p.Add(job)
    p.Add(job)
    p.Add(job)
    p.Add(job)


    time.Sleep(1*time.Second)

    p.Close()

    time.Sleep(5*time.Second)
}
