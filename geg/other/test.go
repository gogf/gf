package main

import (
    "time"
    "gitee.com/johng/gf/g/os/gcache"
    "fmt"
)

func main() {
	c := gcache.New(1000)
	c.Set(1, 1, 0)
	c.Set(2, 2, 0)
	c.Clear()
	fmt.Println(c.Size())
	time.Sleep(time.Second)
}