package main

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmlock"
)

// 内存锁 - TryLock
func main() {
	key := "lock"
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			if gmlock.TryLock(key) {
				glog.Println(i)
				time.Sleep(time.Second)
				gmlock.Unlock(key)
			} else {
				glog.Println(false)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
