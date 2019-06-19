package main

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gmlock"
	"github.com/gogf/gf/g/test/gtest"
	"time"
)

func main() {
	mu := gmlock.NewMutex()
	array := garray.New()
	go func() {
		mu.LockFunc(func() {
			array.Append(1)
			time.Sleep(10000 * time.Millisecond)
		})
	}()
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 10000; i++ {
		go func(i int) {
			time.Sleep(50 * time.Millisecond)
			mu.LockFunc(func() {
				glog.Print(i)
				array.Append(1)
			})
		}(i)
	}
	go func() {
		time.Sleep(60 * time.Millisecond)
		mu.Unlock()
		mu.Unlock()
		mu.Unlock()
	}()

	time.Sleep(20 * time.Millisecond)
	gtest.Assert(array.Len(), 1)
	time.Sleep(50 * time.Millisecond)
	gtest.Assert(array.Len(), 1)
	time.Sleep(50 * time.Millisecond)
	gtest.Assert(array.Len(), 3)
}
