package main

import (
    "gitee.com/johng/gf/g/os/gcache"
    "time"
    "fmt"
)

func main() {
    // 设置LRU淘汰数量
    c := gcache.New(2)

    // 添加10个元素，不过期
    for i := 0; i < 10; i++ {
        c.Set(i, i, 0)
    }
    fmt.Println(c.Size())
    fmt.Println(c.Keys())

    // 读取键名1，保证该键名是优先保留
    fmt.Println(c.Get(1))

    // 等待一定时间后(默认1秒检查一次)，元素会被按照从旧到新的顺序进行淘汰
    time.Sleep(2*time.Second)
    fmt.Println(c.Size())
    fmt.Println(c.Keys())
}