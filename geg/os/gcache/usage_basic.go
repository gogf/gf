package main

import (
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g"
)

func main() {
    // 创建一个缓存对象，当然也可以直接使用gcache包方法
    c := gcache.New()

    // 设置缓存，不过期
    c.Set("k1", "v1", 0)

    // 获取缓存
    g.Dump(c.Get("k1"))

    // 获取缓存大小
    g.Dump(c.Size())

    // 缓存中是否存在指定键名
    g.Dump(c.Contains("k1"))

    // 删除并返回被删除的键值
    g.Dump(c.Remove("k1"))

    // 关闭缓存对象，让GC回收资源
    c.Close()
}