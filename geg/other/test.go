package main

import (
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gwheel"
    "time"
)

func main() {
    v := gtype.NewInt()
    //w := gwheel.New(10, 100*time.Millisecond)
    glog.Println("start")
    for i := 0; i < 10000000; i++ {
        gwheel.AddOnce(time.Second, func() {
           //glog.Println("add")
           v.Add(1)
        })
    }
    glog.Println("end")
    time.Sleep(1100*time.Millisecond)
    glog.Println(v.Val())
    //gwheel.AddSingleton(time.Second, func() {
    //    fmt.Println(time.Now().String())
    //})
    //select { }
}
