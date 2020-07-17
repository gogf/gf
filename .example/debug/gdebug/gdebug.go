package main

import (
	"fmt"
	"github.com/jin502437344/gf/debug/gdebug"
)

func main() {
	gdebug.PrintStack()
	fmt.Println(gdebug.CallerPackage())
	fmt.Println(gdebug.CallerFunction())
}
