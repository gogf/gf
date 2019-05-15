package main

import (
	"os"
	"sync"
)

func main() {
	path := "/Users/john/Temp/test.log"
	os.Remove(path)
	array := make([]*os.File, 1000)
	for i := 0; i < len(array); i++ {
		array[i], _ = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}
	c := make(chan struct{})
	// 62 byte * 10 = 6200 byte
	s := ""
	for i := 0; i < 100; i++ {
		s += "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	s += "\n"
	wg := sync.WaitGroup{}
	wg.Add(1000*len(array))
	for i := 0; i < 1000; i++ {
		go func() {
			<- c
			for i := 0; i < len(array); i++ {
				array[i].WriteString(s)
				wg.Done()
			}
		}()
	}
	close(c)
	wg.Wait()
}