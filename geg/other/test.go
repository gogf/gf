package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gfile"
	"time"
)

func main() {
	go func() {
		go func() {
			fmt.Println("main:", gfile.MainPkgPath())
		}()
	}()
	time.Sleep(time.Second)
}