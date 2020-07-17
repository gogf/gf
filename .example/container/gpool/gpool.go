package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/container/gpool"
)

func main() {
	// 创建一个对象池，过期时间为1000毫秒
	p := gpool.New(1000*time.Millisecond, nil)

	// 从池中取一个对象，返回nil及错误信息
	fmt.Println(p.Get())

	// 丢一个对象到池中
	p.Put(1)

	// 重新从池中取一个对象，返回1
	fmt.Println(p.Get())

	// 等待1秒后重试，发现对象已过期，返回nil及错误信息
	time.Sleep(time.Second)
	fmt.Println(p.Get())
}
