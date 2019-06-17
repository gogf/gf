package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Add(-100)
	wg.Add()
	wg.Wait()
	fmt.Println(1)
}
