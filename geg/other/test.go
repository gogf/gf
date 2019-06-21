package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/os/gmlock"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	mu := gmlock.NewMutex()
	array := garray.New()
	go func() {
		mu.LockFunc(func() {
			array.Append(1)
			time.Sleep(100 * time.Millisecond)
			fmt.Println("====unlock")
		})
	}()
	go func() {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("tryRLock1")
		mu.TryRLockFunc(func() {
			array.Append(1)
			fmt.Println("tryRLock1 success")
		})
	}()
	go func() {
		time.Sleep(150 * time.Millisecond)
		fmt.Println("tryRLock2")
		mu.TryRLockFunc(func() {
			array.Append(1)
			fmt.Println("tryRLock2 success")
		})
	}()
	go func() {
		time.Sleep(150 * time.Millisecond)
		fmt.Println("tryRLock3")
		mu.TryRLockFunc(func() {
			array.Append(1)
			fmt.Println("tryRLock3 success")
		})
	}()
	time.Sleep(50 * time.Millisecond)
	gtest.Assert(array.Len(), 1)
	time.Sleep(50 * time.Millisecond)
	gtest.Assert(array.Len(), 1)
	time.Sleep(150 * time.Millisecond)
	fmt.Println("====array len:", array.Len())
	gtest.Assert(array.Len(), 3)
}
