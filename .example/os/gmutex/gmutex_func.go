package main

import (
	"time"

	"github.com/gogf/gf/v2/os/glog"

	"github.com/gogf/gf/v2/os/gmutex"
)

func main() {
	mu := gmutex.New()
	go mu.LockFunc(func() {
		glog.Print("lock func1")
		time.Sleep(1 * time.Second)
	})
	time.Sleep(time.Millisecond)
	go mu.LockFunc(func() {
		glog.Print("lock func2")
	})
	time.Sleep(2 * time.Second)
}
