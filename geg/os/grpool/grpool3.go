package main

import (
<<<<<<< HEAD
    "fmt"
    "sync"
)

func main() {
    wg := sync.WaitGroup{}
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(v int){
            fmt.Println(v)
            wg.Done()
        }(i)
    }
    wg.Wait()
=======
	"fmt"
	"github.com/gogf/gf/g/os/grpool"
	"sync"
)

func main() {
	p  := grpool.New(1)
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
>>>>>>> upstream/master
}
