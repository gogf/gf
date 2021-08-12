package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gogf/gf/os/gmlock"
)

// 测试Locker是否会产生死锁
func main() {
	var (
		l      = gmlock.New()
		wg     = sync.WaitGroup{}
		key    = "test"
		event  = make(chan int)
		number = 100000
	)
	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			l.Lock(key)
			//fmt.Println("get lock")
			l.Unlock(key)
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			l.RLock(key)
			//fmt.Println("get rlock")
			l.RUnlock(key)
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if l.TryLock(key) {
				//fmt.Println("get lock")
				l.Unlock(key)
			}
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if l.TryRLock(key) {
				//fmt.Println("get rlock")
				l.RUnlock(key)
			}
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if l.TryLock(key) {
				// 模拟业务逻辑的随机处理间隔
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				l.Unlock(key)
			}
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if l.TryRLock(key) {
				// 模拟业务逻辑的随机处理间隔
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				l.RUnlock(key)
			}
			wg.Done()
		}()
	}
	// 使用chan作为事件发送测试指令，让所有的goroutine同时执行
	close(event)
	wg.Wait()

	fmt.Println("done!")
}
