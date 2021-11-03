package main

import (
	"fmt"
	"sync"

	"github.com/gogf/gf/v2/os/grpool"
)

func main() {
	p := grpool.New(1)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		v := i
		p.Add(func() {
			fmt.Println(v)
			wg.Done()
		})
	}
	wg.Wait()
}
