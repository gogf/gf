package main

import (
	"fmt"

	"github.com/jin502437344/gf/os/gcache"
)

func main() {
	// 创建一个缓存对象，当然也可以直接使用gcache包方法
	c := gcache.New()

	// 设置缓存，不过期
	c.Set("k1", "v1", 0)

	// 获取缓存
	fmt.Println(c.Get("k1"))

	// 获取缓存大小
	fmt.Println(c.Size())

	// 缓存中是否存在指定键名
	fmt.Println(c.Contains("k1"))

	// 删除并返回被删除的键值
	fmt.Println(c.Remove("k1"))

	// 关闭缓存对象，让GC回收资源
	c.Close()
}
