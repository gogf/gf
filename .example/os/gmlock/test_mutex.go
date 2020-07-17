package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/jin502437344/gf/os/gmlock"
)

// 测试是否会产生死锁
func main() {
	mu := gmlock.NewMutex()
	wg := sync.WaitGroup{}
	event := make(chan int)
	number := 100000

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			mu.Lock()
			//fmt.Println("get lock")
			mu.Unlock()
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			mu.RLock()
			//fmt.Println("get rlock")
			mu.RUnlock()
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if mu.TryLock() {
				//fmt.Println("get lock")
				mu.Unlock()
			}
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if mu.TryRLock() {
				//fmt.Println("get rlock")
				mu.RUnlock()
			}
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if mu.TryLock() {
				// 模拟业务逻辑的随机处理间隔
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				mu.Unlock()
			}
			wg.Done()
		}()
	}

	for i := 0; i < number; i++ {
		wg.Add(1)
		go func() {
			<-event
			if mu.TryRLock() {
				// 模拟业务逻辑的随机处理间隔
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				mu.RUnlock()
			}
			wg.Done()
		}()
	}
	// 使用chan作为事件发送测试指令，让所有的goroutine同时执行
	close(event)
	wg.Wait()

	fmt.Println("done!")
}
