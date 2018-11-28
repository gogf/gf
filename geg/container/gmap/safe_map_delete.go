package main

import (
    "sync"
    "time"
    "fmt"
)

// 验证 map 的delete方法是否并发安全
func main() {
    // 创建一个初始化的map
    m := make(map[int]int)
    for i := 0; i < 10000; i++ {
        m[i] = i
    }

    fmt.Println("map size:", len(m))

    wg := sync.WaitGroup{}
    ev := make(chan struct{}, 0)

    // 创建10个并发的goroutine，使用ev控制并发开始事件，更容易模拟data race
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            <- ev
            fmt.Println("start")
            for i := 0; i < 10000; i++ {
                delete(m, i)
            }
            wg.Done()
        }()
    }

    time.Sleep(time.Second)

    close(ev)
    wg.Wait()
}
