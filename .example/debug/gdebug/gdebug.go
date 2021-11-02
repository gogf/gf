package main

import (
	"fmt"
	"github.com/gogf/gf/v2/debug/gdebug"
)

func main() {
	gdebug.PrintStack()
	fmt.Println(gdebug.CallerPackage())
	fmt.Println(gdebug.CallerFunction())
}
