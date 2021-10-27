package main

import (
	"fmt"
	"github.com/gogf/gf/v2/container/glist"
)

func main() {
	l := glist.New()
	// Push
	l.PushBack(1)
	l.PushBack(2)
	e := l.PushFront(0)
	// Insert
	l.InsertBefore(e, -1)
	l.InsertAfter(e, "a")
	fmt.Println(l)
	// Pop
	fmt.Println(l.PopFront())
	fmt.Println(l.PopBack())
	fmt.Println(l)
}
