package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gmlock"
)

// 内存锁 - 手动Unlock与计时Unlock冲突校验
func main() {
	key := "key"

	// 第一次锁带时间
	gmlock.Lock(key)
	glog.Println("lock1")
	// 这个时候上一次的计时解锁已失效
	gmlock.Unlock(key)
	glog.Println("unlock1")

	fmt.Println()

	// 第二次锁，不带时间，且在执行过程中钱一个Lock的定时解锁生效
	gmlock.Lock(key)
	glog.Println("lock2")
	go func() {
		// 正常情况下3秒后才能执行这句
		gmlock.Lock(key)
		glog.Println("lock by goroutine")
	}()
	time.Sleep(3 * time.Second)
	// 这时再解锁
	gmlock.Unlock(key)
	// 注意3秒之后才会执行这一句
	glog.Println("unlock2")

	select {}
}
