package main

import (
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g"
    "time"
)

func main() {
    // 当键名不存在时写入，设置过期时间1000毫秒
    gcache.SetIfNotExist("k1", "k1", 1000)

    // 打印当前的键名列表
    g.Dump(gcache.Keys())

    // 打印当前的键名列表 []string 类型
    g.Dump(gcache.KeyStrings())

    // 获取指定键值，如果不存在时写入，并返回键值
    g.Dump(gcache.GetOrSet("k2", "v2", 0))

    // 打印当前的键值列表
    g.Dump(gcache.Values())

    // 等待1秒，以便k1:v1自动过期
    time.Sleep(time.Second)

    // 再次打印当前的键值列表，发现k1:v1已经过期，只剩下k2:v2
    g.Dump(gcache.Values())
}