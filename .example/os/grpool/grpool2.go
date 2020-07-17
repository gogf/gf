package main

import (
	"fmt"
	"sync"

	"github.com/jin502437344/gf/os/grpool"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		v := i
		grpool.Add(func() {
			fmt.Println(v)
			wg.Done()
		})
	}
	wg.Wait()
}
