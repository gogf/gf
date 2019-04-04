package main

import (
	"fmt"
	"github.com/gogf/gf/g/container/gring"
)

func main() {
	r := gring.New(3)
	r.Put(1)
	r.Put(2)
	fmt.Println(r.Val())
}