package main

import (
<<<<<<< HEAD
    "fmt"
    "github.com/theckman/go-flock"
    "time"
)

func main() {
    l := flock.NewFlock("/tmp/go-lock.lock")
    l.Lock()
    fmt.Printf("lock 1")
    l.Lock()
    fmt.Printf("lock 1")

    time.Sleep(time.Hour)
=======
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
>>>>>>> upstream/master
}
