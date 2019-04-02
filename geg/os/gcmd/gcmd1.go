package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gcmd"
)

func help() {
	fmt.Println("This is help.")
}

func test() {
	fmt.Println("This is test.")
}

func main() {
	gcmd.BindHandle("help", help)
	gcmd.BindHandle("test", test)
	gcmd.AutoRun()
}
