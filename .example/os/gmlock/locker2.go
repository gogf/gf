package main

import (
	"sync"

	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gmlock"
)

// 内存锁 - 给定过期时间
func main() {
	key := "lock"
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			gmlock.Lock(key, 1000)
			glog.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
