package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/g/os/grpool"
)

func main() {
	grpool.Add(func() {

	})
	fmt.Println(grpool.Size())
	fmt.Println(grpool.Jobs())
	time.Sleep(time.Second)
	fmt.Println(grpool.Size())
	fmt.Println(grpool.Jobs())

}
