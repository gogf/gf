package main

import (
	"github.com/gogf/gf/frame/g"
	"sync"
)

func main() {
	var (
		wg = sync.WaitGroup{}
		ch = make(chan struct{})
	)
	wg.Add(3000)
	for i := 0; i < 3000; i++ {
		go func() {
			<-ch
			g.Log().Println("abcdefghijklmnopqrstuvwxyz1234567890")
			wg.Done()
		}()
	}
	close(ch)
	wg.Wait()
}
