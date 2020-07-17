package main

import (
	"sync"

	"github.com/jin502437344/gf/os/glog"
)

func main() {
	wg := sync.WaitGroup{}
	c := make(chan struct{})
	wg.Add(3000)
	for i := 0; i < 3000; i++ {
		go func() {
			<-c
			glog.Println("abcdefghijklmnopqrstuvwxyz1234567890")
			wg.Done()
		}()
	}
	close(c)
	wg.Wait()
}
