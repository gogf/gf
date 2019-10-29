package main

import (
	"fmt"
	"github.com/gogf/gf/container/glist"
)

func main() {
	l := glist.New()
	// Push
	l.PushBack(1)
	l.PushBack(2)
	e0 := l.PushFront(0)
	// Insert
	l.InsertBefore(e0, -1)
	l.InsertAfter(e0, "a")
	fmt.Println(l)
	// Pop
	fmt.Println(l.PopFront())
	fmt.Println(l.PopBack())
	fmt.Println(l)
}
