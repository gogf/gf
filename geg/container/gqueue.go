package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/gqueue"
    "time"
)

func main() {
    t := gtime.Microsecond()
    q := gqueue.NewInterfaceQueue()
    fmt.Println("queue creation costs(μs):", gtime.Microsecond() - t)

    // 每隔2秒异步打印出当前队列的大小
    gtime.SetInterval(2*time.Second, func() bool {
        fmt.Println("queue size:", q.Size())
        return true
    })

    // push10条数据
    for i := 0; i < 10; i++ {
        q.Push(i)
        fmt.Println("push:", i)
    }

    // 每隔1秒pop1条数据
    for {
        time.Sleep(time.Second)
        fmt.Println(" pop:", q.Pop())
    }
}
