package main

import (
	"fmt"
	"sync"
)

func main() {
	m := sync.RWMutex{}
	m.Lock()
	fmt.Println(m)
}
