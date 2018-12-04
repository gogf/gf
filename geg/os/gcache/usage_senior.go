package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gcache"
    "time"
)

func main() {
    // 当键名不存在时写入，设置过期时间1000毫秒
    gcache.SetIfNotExist("k1", "v1", 1000)

    // 打印当前的键名列表
    fmt.Println(gcache.Keys())

    // 打印当前的键值列表
    fmt.Println(gcache.Values())

    // 获取指定键值，如果不存在时写入，并返回键值
    fmt.Println(gcache.GetOrSet("k2", "v2", 0))

    // 打印当前的键值对
    fmt.Println(gcache.Data())

    // 等待1秒，以便k1:v1自动过期
    time.Sleep(time.Second)

    // 再次打印当前的键值对，发现k1:v1已经过期，只剩下k2:v2
    fmt.Println(gcache.Data())
}