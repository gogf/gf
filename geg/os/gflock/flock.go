package main

import (
	"fmt"
	"github.com/gogf/gf/third/github.com/theckman/go-flock"
	"time"
)

func main() {
	l := flock.NewFlock("/tmp/go-lock.lock")
	l.Lock()
	fmt.Printf("lock 1")
	l.Lock()
	fmt.Printf("lock 1")

	time.Sleep(time.Hour)
}
