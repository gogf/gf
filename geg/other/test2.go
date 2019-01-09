package main

import (
    "container/list"
    "gitee.com/johng/gf/g/os/glog"
    "time"
)

func main(){
    list := list.New()
    glog.Println("start1")
    for i := 0; i < 10000000; i++ {
        list.PushBack(i)
    }
    glog.Println("end1")

    glog.Println("start2")
    for e := list.Front(); e != nil; e = e.Next() {
        time.Sleep(25*time.Nanosecond)
    }
    glog.Println("end2")
}
