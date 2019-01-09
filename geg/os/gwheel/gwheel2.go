package main

import (
    "fmt"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gwheel"
    "time"
)

func main() {
    v := gtype.NewInt()
    w := gwheel.New(10, 10*time.Millisecond)
    fmt.Println("start:", time.Now())
    for i := 0; i < 1000000; i++ {
        w.AddOnce(time.Second, func() {
            v.Add(1)
        })
    }
    fmt.Println("end  :", time.Now())
    time.Sleep(3020*time.Millisecond)
    fmt.Println(v.Val(), time.Now())

    //gwheel.AddSingleton(time.Second, func() {
    //    fmt.Println(time.Now().String())
    //})
    //select { }
}
