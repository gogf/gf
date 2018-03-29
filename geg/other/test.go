package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gidgen"
)

func main() {
    g := gidgen.New(2)
    for i := 0; i < 11; i++ {
        fmt.Println(g.Int())
    }
    g.Close()
    fmt.Println(g.Uint())
    //events2 := make(chan int, 100)
    //go func() {
    //    for{
    //        v := <- events1
    //        fmt.Println(v)
    //    }
    //
    //}()

    //go func() {
    //    time.Sleep(2*time.Second)
    //    events1 <- 1
    //    events2 <- 2
    //    time.Sleep(2*time.Second)
    //    close(events1)
    //    close(events2)
    //    events1 <- 1
    //    events2 <- 2
    //}()
    //
    //select {
    //
    //}
}