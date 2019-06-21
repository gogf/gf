package main

import (
	"time"

	"github.com/gogf/gf/g/os/gmutex"

	"github.com/gogf/gf/g/os/glog"

	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	mu := gmutex.New()
	array := garray.New()
	glog.Println("step0")
	go func() {
		mu.LockFunc(func() {
			array.Append(1)
			time.Sleep(200 * time.Millisecond)
			glog.Println("unlocked")
		})
	}()
	go func() {
		time.Sleep(150 * time.Millisecond)
		mu.TryRLockFunc(func() {
			array.Append(1)
			glog.Println("add1")
		})
	}()
	sum := 1000
	for index := 1; index < sum; index++ {
		go func(i int) {
			time.Sleep(300 * time.Millisecond)
			//fmt.Println(i*10, mu.IsLocked())
			if r := mu.TryRLockFunc(func() {
				array.Append(1)
				time.Sleep(200 * time.Millisecond)
			}); !r {
				glog.Println(i, r)
			}
		}(index)
	}
	glog.Println("step1")
	time.Sleep(100 * time.Millisecond)
	glog.Println("step2")
	gtest.Assert(array.Len(), 1)
	time.Sleep(100 * time.Millisecond)
	glog.Println("step3")
	gtest.Assert(array.Len(), 1)
	time.Sleep(1000 * time.Millisecond)
	gtest.Assert(array.Len(), sum)

}
