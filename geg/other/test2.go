package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gmlock"
	"time"
)

func main() {
	key   := "test3"
	gmlock.Lock(key, 200*time.Millisecond)
	fmt.Println("TryLock:", gmlock.TryLock(key))
	fmt.Println("TryLock:", gmlock.TryLock(key))
}
