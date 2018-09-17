package main

import (
    "gitee.com/johng/gf/g/os/gcache"
    "time"
    "fmt"
)

func main() {
    // 设置LRU淘汰数量
    gcache.SetCap(2)

    // 10个元素
    for i := 0; i < 10; i++ {
        gcache.Set(i, i, 0)
    }
    fmt.Println(gcache.Size())
    fmt.Println(gcache.Keys())

    // 等待一定时间后(默认10秒检查一次)，元素会被按照从旧到新的顺序进行淘汰
    time.Sleep(11*time.Second)
    fmt.Println(gcache.Size())
    fmt.Println(gcache.Keys())
}