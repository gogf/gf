package main

import (
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmutex"
)

func main() {
	mu := gmutex.New()
	for i := 0; i < 10; i++ {
		go func(n int) {
			mu.Lock()
			defer mu.Unlock()
			glog.Println("Lock:", n)
			time.Sleep(time.Second)
		}(i)
	}
	for i := 0; i < 10; i++ {
		go func(n int) {
			mu.RLock()
			defer mu.RUnlock()
			glog.Println("RLock:", n)
			time.Sleep(time.Second)
		}(i)
	}
	time.Sleep(11 * time.Second)
}
