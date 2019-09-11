package main

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
)

func main() {
	fmt.Println(gdebug.CallerPackage())
	fmt.Println(gdebug.CallerFunction())
}
